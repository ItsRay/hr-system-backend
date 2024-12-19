package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"hr-system/internal/common"
	employee_domain "hr-system/internal/employees/domain"
	mocks_employee_repo "hr-system/internal/employees/repo/mocks"
	mocks_leave_cache "hr-system/internal/leaves/cache/mocks"
	"hr-system/internal/leaves/domain"
	mocks_leave_repo "hr-system/internal/leaves/repo/mocks"
)

func genFakeLeave() domain.Leave {
	startDate := time.Now().Truncate(time.Second)
	endDate := startDate.Add(time.Hour * 24)

	leave := domain.Leave{
		ID:                1, // You can assign auto-incremented ID if necessary
		EmployeeID:        3,
		Type:              domain.LeaveTypeAnnual,
		StartDate:         startDate,
		EndDate:           endDate,
		Reason:            "Vacation",
		Status:            domain.ReviewStatusReviewing,
		CurrentReviewerID: common.GetPtr(2),
		Reviews: []domain.LeaveReview{
			{
				ID:         1,
				LeaveID:    1,
				ReviewerID: 2,
				Status:     domain.ReviewStatusReviewing,
			},
		},
	}

	return leave
}

func TestCreateLeave(t *testing.T) {
	mockLeaveRepo := mocks_leave_repo.NewLeaveRepo(t)
	mockEmployeeRepo := mocks_employee_repo.NewEmployeeRepo(t)
	mockLeaveCache := mocks_leave_cache.NewLeaveCache(t)
	logger := common.NewLogger()

	service := NewLeaveService(logger, mockLeaveRepo, mockEmployeeRepo, mockLeaveCache)

	ctx := context.Background()
	leave := genFakeLeave()

	mockEmployeeRepo.On("GetEmployeeByID", ctx, leave.EmployeeID).
		Return(employee_domain.Employee{ID: 1, ManagerID: nil}, nil).Once()
	mockLeaveRepo.On("CreateLeave", ctx, &leave).Return(nil).Once()
	mockLeaveCache.On("DelLeavesFromCache", ctx, mock.Anything).Return(nil).Once()
	mockLeaveCache.On("SetLeaveToCache", ctx, &leave).Return(nil).Once()

	createdLeave, err := service.CreateLeave(ctx, &leave)
	assert.NoError(t, err)
	assert.Equal(t, leave, createdLeave)
}

func TestReviewLeave(t *testing.T) {
	mockLeaveRepo := mocks_leave_repo.NewLeaveRepo(t)
	mockEmployeeRepo := mocks_employee_repo.NewEmployeeRepo(t)
	mockLeaveCache := mocks_leave_cache.NewLeaveCache(t)
	logger := common.NewLogger()

	service := NewLeaveService(logger, mockLeaveRepo, mockEmployeeRepo, mockLeaveCache)

	ctx := context.Background()

	leave := genFakeLeave()
	reviewerID := leave.Reviews[0].ReviewerID

	mockLeaveRepo.On("GetLeaveByID", ctx, leave.ID).Return(leave, nil).Once()
	mockEmployeeRepo.On("GetEmployeeByID", ctx, reviewerID).
		Return(employee_domain.Employee{ID: 1,
			Positions: []employee_domain.Position{
				{
					ManagerLevel: 5,
				},
			}}, nil).
		Once()
	mockLeaveRepo.On("UpdateLeaveAndReviews", ctx, mock.Anything, mock.Anything).Return(nil).Once()
	mockLeaveCache.On("DelLeaveFromCache", ctx, leave.ID).Return(nil).Once()
	mockLeaveCache.On("DelLeavesFromCache", ctx, mock.Anything).Return(nil).Twice()

	err := service.ReviewLeave(ctx, leave.ID, reviewerID, domain.ReviewStatusApproved, "")
	assert.NoError(t, err)
}
