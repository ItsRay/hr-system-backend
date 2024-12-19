package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"hr-system/internal/leaves/domain"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&domain.Leave{}, &domain.LeaveReview{}); err != nil {
		return nil, err
	}
	return db, nil
}

func TestCreateLeave(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := &leaveRepo{db: db}
	leave := &domain.Leave{EmployeeID: 1, Reason: "Vacation"}

	err = repo.CreateLeave(context.Background(), leave)
	assert.NoError(t, err)
	assert.NotZero(t, leave.ID)
}

func TestGetLeaveByID(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := &leaveRepo{db: db}
	leave := &domain.Leave{EmployeeID: 1, Reason: "Vacation"}
	err = repo.CreateLeave(context.Background(), leave)
	assert.NoError(t, err)

	fetchedLeave, err := repo.GetLeaveByID(context.Background(), leave.ID)
	assert.NoError(t, err)
	assert.Equal(t, leave.ID, fetchedLeave.ID)
	assert.Equal(t, leave.EmployeeID, fetchedLeave.EmployeeID)
	assert.Equal(t, leave.Reason, fetchedLeave.Reason)
}

func TestUpdateLeaveAndReviews(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := &leaveRepo{db: db}
	leave := &domain.Leave{EmployeeID: 2, Reason: "Vacation",
		Reviews: []domain.LeaveReview{
			{
				ReviewerID: 1,
				Status:     domain.ReviewStatusReviewing,
			},
		},
	}
	err = repo.CreateLeave(context.Background(), leave)
	assert.NoError(t, err)

	leave.Status = domain.ReviewStatusApproved
	reviews := []domain.LeaveReview{leave.Reviews[0]}
	reviews[0].Status = domain.ReviewStatusApproved

	err = repo.UpdateLeaveAndReviews(context.Background(), leave, reviews)
	assert.NoError(t, err)

	fetchedLeave, err := repo.GetLeaveByID(context.Background(), leave.ID)
	assert.NoError(t, err)
	assert.Equal(t, leave.Status, fetchedLeave.Status)
	assert.Equal(t, 1, len(fetchedLeave.Reviews))
	assert.Equal(t, reviews[0].Status, fetchedLeave.Reviews[0].Status)
}

func TestGetLeaves(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := &leaveRepo{db: db}
	leave1 := &domain.Leave{EmployeeID: 1, Reason: "Vacation"}
	leave2 := &domain.Leave{EmployeeID: 2, Reason: "Sick Leave"}

	err = repo.CreateLeave(context.Background(), leave1)
	assert.NoError(t, err)
	err = repo.CreateLeave(context.Background(), leave2)
	assert.NoError(t, err)

	query := domain.LeavesQuery{EmployeeID: &leave1.EmployeeID}
	leaves, err := repo.GetLeaves(context.Background(), query)
	assert.NoError(t, err)
	assert.Len(t, leaves, 1)
	assert.Equal(t, leave1.EmployeeID, leaves[0].EmployeeID)
	assert.Equal(t, leave1.Reason, leaves[0].Reason)
}
