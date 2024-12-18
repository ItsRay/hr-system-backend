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
	GetLeaves(ctx context.Context, query domain.LeaveQuery) ([]domain.Leave, error)
	ReviewLeave(ctx context.Context, leaveID, reviewerID int, decision domain.ReviewStatus, comment string) error
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
		return domain.Leave{}, err
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
		return fmt.Errorf("invalid decision: %s", decision)
	}

	leave, err := s.leaveRepo.GetLeaveByID(ctx, leaveID)
	if err != nil {
		if errors.Is(err, common.ErrResourceNotFound) {
			return common.ErrResourceNotFound
		}
		return fmt.Errorf("failed to retrieve leave: %w", err)
	}
	if leave.Status != domain.ReviewStatusReviewing {
		return fmt.Errorf("leave is not in reviewing status")
	}

	// check reviewer permission
	if leave.CurrentReviewerID == nil || *leave.CurrentReviewerID != reviewerID {
		return fmt.Errorf("reviewer does not have permission to review this leave")
	}

	// update review comment & status
	now := time.Now()
	leave.Reviews[len(leave.Reviews)-1].Status = decision
	leave.Reviews[len(leave.Reviews)-1].Comment = comment
	leave.Reviews[len(leave.Reviews)-1].ReviewedAt = &now

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
			leave.Reviews = append(leave.Reviews, domain.LeaveReview{
				ReviewerID: *reviewer.ManagerID,
				Status:     domain.ReviewStatusReviewing,
			})
			leave.CurrentReviewerID = reviewer.ManagerID
		} else {
			// leave approved
			leave.CurrentReviewerID = nil
			leave.Status = domain.ReviewStatusApproved
		}
	} else {
		// rejected
		leave.Status = domain.ReviewStatusRejected
		leave.CurrentReviewerID = nil
		leave.Reviews[len(leave.Reviews)-1].Comment = comment
	}

	err = s.leaveRepo.UpdateLeave(ctx, &leave)
	if err != nil {
		return fmt.Errorf("failed to update leave review: %w", err)
	}

	return nil
}

func (s *leaveService) GetLeaves(ctx context.Context, query domain.LeaveQuery) ([]domain.Leave, error) {
	// TODO: tell it's 4xx or 5xx error
	if query.EmployeeID == nil && query.CurrentReviewerID == nil {
		return nil, fmt.Errorf("employee ID or current reviewer ID must be provided")
	}

	leaves, err := s.leaveRepo.GetLeaves(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaves: %w", err)
	}

	return leaves, nil
}