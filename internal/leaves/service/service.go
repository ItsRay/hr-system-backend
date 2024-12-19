package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"

	"hr-system/internal/common"
	employee_domain "hr-system/internal/employees/domain"
	employee_repo "hr-system/internal/employees/repo"
	"hr-system/internal/leaves/cache"
	"hr-system/internal/leaves/domain"
	"hr-system/internal/leaves/repo"
)

type LeaveService interface {
	CreateLeave(ctx context.Context, leave *domain.Leave) (domain.Leave, error)
	GetLeaves(ctx context.Context, query domain.LeavesQuery) ([]domain.Leave, error)
	ReviewLeave(ctx context.Context, leaveID, reviewerID int, decision domain.ReviewStatus, comment string) error
	GetLeaveByID(ctx context.Context, id int) (domain.Leave, error)
}

type leaveService struct {
	leaveRepo    repo.LeaveRepo
	leaveCache   cache.LeaveCache
	employeeRepo employee_repo.EmployeeRepo
	logger       *common.Logger
	validate     *validator.Validate
}

func NewLeaveService(logger *common.Logger, leaveRepo repo.LeaveRepo, employeeRepo employee_repo.EmployeeRepo,
	leaveCache cache.LeaveCache) LeaveService {
	return &leaveService{
		leaveRepo:    leaveRepo,
		employeeRepo: employeeRepo,
		leaveCache:   leaveCache,
		logger:       logger,
		validate:     validator.New(),
	}
}

func (s *leaveService) validateCreateLeave(leave *domain.Leave) error {
	if err := s.validate.Struct(leave); err != nil {
		return fmt.Errorf("failed to validate leave: %w", err)
	}

	if leave.StartDate.After(leave.EndDate) {
		return fmt.Errorf("start date must be before end date")
	}

	return nil
}

func (s *leaveService) CreateLeave(ctx context.Context, leave *domain.Leave) (domain.Leave, error) {
	if err := s.validateCreateLeave(leave); err != nil {
		return domain.Leave{}, fmt.Errorf("%w, detail: %s", common.ErrInvalidInput, err)
	}

	employee, err := s.employeeRepo.GetEmployeeByID(ctx, leave.EmployeeID)
	if err != nil {
		if errors.Is(err, common.ErrResourceNotFound) {
			return domain.Leave{}, common.ErrResourceNotFound
		}
		return domain.Leave{}, fmt.Errorf("failed to get manager IDs: %w", err)
	}

	// status
	leave.Status = domain.ReviewStatusReviewing
	if employee.ManagerID == nil {
		leave.Status = domain.ReviewStatusApproved
	}

	// currentReviewerID
	leave.CurrentReviewerID = employee.ManagerID

	// reviews
	if employee.ManagerID != nil {
		leave.Reviews = []domain.LeaveReview{
			{
				ReviewerID: *employee.ManagerID,
				Status:     domain.ReviewStatusReviewing,
			},
		}
	}

	if err := s.leaveRepo.CreateLeave(ctx, leave); err != nil {
		return domain.Leave{}, fmt.Errorf("failed to create leave: %w", err)
	}

	// delete cache of this employee
	if err := s.leaveCache.DelLeavesFromCache(ctx, domain.LeavesQuery{EmployeeID: &leave.EmployeeID}); err != nil {
		s.logger.Errorf("failed to delete employee %d cache, cause: %s", leave.EmployeeID, err)
	}
	// delete cache of this reviewer
	if leave.CurrentReviewerID != nil {
		err := s.leaveCache.DelLeavesFromCache(ctx, domain.LeavesQuery{CurrentReviewerID: leave.CurrentReviewerID})
		if err != nil {
			s.logger.Errorf("failed to delete reviewer %d cache, cause: %s", leave.CurrentReviewerID, err)
		}
	}

	if err := s.leaveCache.SetLeaveToCache(ctx, leave); err != nil {
		s.logger.Warnf("Failed to cache leave data: %v", err)
	}

	return *leave, nil
}

// check if the leave needs to be reviewed by the next reviewer
func needNextReviewer(leave *domain.Leave, approver *employee_domain.Employee) bool {
	days := int(leave.EndDate.Sub(leave.StartDate).Hours() / 24)

	var needManagerLevel int
	if days <= 5 {
		needManagerLevel = 0
	} else if days <= 10 {
		needManagerLevel = 3
	} else {
		needManagerLevel = 5
	}

	return approver.Positions[0].ManagerLevel < needManagerLevel
}

func (s *leaveService) ReviewLeave(ctx context.Context, leaveID int, reviewerID int, decision domain.ReviewStatus,
	comment string) error {
	if decision != domain.ReviewStatusApproved && decision != domain.ReviewStatusRejected {
		return fmt.Errorf("%w, invalid decision: %s", common.ErrInvalidInput, decision)
	}

	leave, err := s.leaveRepo.GetLeaveByID(ctx, leaveID)
	if err != nil {
		if errors.Is(err, common.ErrResourceNotFound) {
			return common.ErrResourceNotFound
		}
		return fmt.Errorf("failed to retrieve leave: %w", err)
	}
	if leave.Status != domain.ReviewStatusReviewing {
		return fmt.Errorf("%w, leave is not in reviewing status", common.ErrStatusConflict)
	}

	// check reviewer permission
	if leave.CurrentReviewerID == nil || *leave.CurrentReviewerID != reviewerID {
		return fmt.Errorf("%w, it's not waiting for this reviewer to review", common.ErrStatusConflict)
	}

	// update review comment & status
	now := time.Now()
	if len(leave.Reviews) == 0 {
		return fmt.Errorf("unexpected error: no review found")
	}
	updateReviews := []domain.LeaveReview{
		leave.Reviews[len(leave.Reviews)-1],
	}
	updateReviews[0].Comment = comment
	updateReviews[0].ReviewedAt = &now
	updateReviews[0].Status = decision

	// update leave
	if decision == domain.ReviewStatusApproved {
		// approved
		reviewer, err := s.employeeRepo.GetEmployeeByID(ctx, reviewerID)
		if err != nil {
			if errors.Is(err, common.ErrResourceNotFound) {
				return common.ErrResourceNotFound
			}
			return fmt.Errorf("failed to get manager IDs: %w", err)
		}
		if needNextReviewer(&leave, &reviewer) {
			// pass to next reviewer
			if reviewer.ManagerID == nil {
				return fmt.Errorf("unexpected error: reviewer does not have a manager")
			}
			updateReviews = append(updateReviews, domain.LeaveReview{
				LeaveID:    leaveID,
				ReviewerID: *reviewer.ManagerID,
				Status:     domain.ReviewStatusReviewing,
			})
			leave.CurrentReviewerID = reviewer.ManagerID
		} else {
			// leave approved
			leave.Status = domain.ReviewStatusApproved
			leave.CurrentReviewerID = nil
		}
	} else {
		// rejected
		leave.Status = domain.ReviewStatusRejected
		leave.CurrentReviewerID = nil
	}

	err = s.leaveRepo.UpdateLeaveAndReviews(ctx, &leave, updateReviews)
	if err != nil {
		return fmt.Errorf("failed to update leave review: %w", err)
	}

	// delete cache of this leave
	if err := s.leaveCache.DelLeaveFromCache(ctx, leaveID); err != nil {
		s.logger.Errorf("failed to delete leave %d cache, cause: %s", leaveID, err)
	}
	// delete cache of this employee
	if err := s.leaveCache.DelLeavesFromCache(ctx, domain.LeavesQuery{EmployeeID: &leave.EmployeeID}); err != nil {
		s.logger.Errorf("failed to delete employee %d cache, cause: %s", leave.EmployeeID, err)
	}
	// delete cache of old reviewer
	if err := s.leaveCache.DelLeavesFromCache(ctx, domain.LeavesQuery{CurrentReviewerID: &reviewerID}); err != nil {
		s.logger.Errorf("failed to delete reviewer %d cache, cause: %s", reviewerID, err)
	}
	// delete cache of new reviewer
	if leave.CurrentReviewerID != nil && *leave.CurrentReviewerID != reviewerID {
		err = s.leaveCache.DelLeavesFromCache(ctx, domain.LeavesQuery{CurrentReviewerID: leave.CurrentReviewerID})
		if err != nil {
			s.logger.Errorf("failed to delete reviewer %d cache, cause: %s", *leave.CurrentReviewerID, err)
		}
	}

	return nil
}

func (s *leaveService) GetLeaves(ctx context.Context, query domain.LeavesQuery) ([]domain.Leave, error) {
	if query.EmployeeID == nil && query.CurrentReviewerID == nil {
		return nil, fmt.Errorf("%w, employee ID or current reviewer ID must be provided", common.ErrInvalidInput)
	}
	if query.EmployeeID != nil && query.CurrentReviewerID != nil {
		// for index and cache
		return nil, fmt.Errorf("%w, only one of employee ID or current reviewer ID can be provided", common.ErrInvalidInput)
	}

	// get from cache
	leaves, err := s.leaveCache.GetLeavesFromCache(ctx, query)
	if err == nil {
		s.logger.Infof("[Cache Hit] leaves query: %+v", query)
		return leaves, nil
	}
	if err != nil && !errors.Is(err, common.ErrResourceNotFound) {
		s.logger.Warnf("failed to get leaves from cache, cause: %s", err)
	}

	// get from repo
	leaves, err = s.leaveRepo.GetLeaves(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaves: %w", err)
	}

	// set cache
	if err := s.leaveCache.SetLeavesToCache(ctx, query, leaves); err != nil {
		s.logger.Errorf("failed to cache leaves data: %v", err)
	}

	return leaves, nil
}

func (s *leaveService) GetLeaveByID(ctx context.Context, id int) (domain.Leave, error) {
	leave, err := s.leaveCache.GetLeaveFromCache(ctx, id)
	if err == nil {
		s.logger.Infof("[Cache Hit] leave id: %d", id)
		return leave, nil
	}
	if err != nil && !errors.Is(err, common.ErrResourceNotFound) {
		s.logger.Warnf("failed to get leave[%d] from cache, cause: %s", id, err)
	}

	leave, err = s.leaveRepo.GetLeaveByID(ctx, id)
	if err != nil {
		if errors.Is(err, common.ErrResourceNotFound) {
			return domain.Leave{}, common.ErrResourceNotFound
		}
		return domain.Leave{}, fmt.Errorf("failed to get leave: %w", err)
	}

	if err := s.leaveCache.SetLeaveToCache(ctx, &leave); err != nil {
		s.logger.Warnf("failed to cache leave data: %v", err)
	}

	return leave, nil
}
