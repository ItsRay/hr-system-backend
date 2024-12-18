package repo

import (
	"context"
	"fmt"
	"time"

	employee_repo "hr-system/internal/employees/repo"
	"hr-system/internal/leaves/domain"
)

func (r *leaveRepo) SeedLeaveData(ctx context.Context, employeeRepo employee_repo.EmployeeRepo) error {
	employees, _, err := employeeRepo.GetEmployees(ctx, 1, 1000)
	if err != nil {
		return fmt.Errorf("failed to fetch employees: %w", err)
	}

	leaves := []domain.Leave{
		{
			EmployeeID: employees[0].ID, // Alice Johnson
			Type:       domain.LeaveTypeAnnual,
			StartDate:  time.Now().AddDate(0, 0, -5), // 5 days ago
			EndDate:    time.Now().AddDate(0, 0, -3), // 3 days ago
			Reason:     "Vacation",
			Status:     domain.ReviewStatusApproved,
		},
		{
			EmployeeID: employees[1].ID,
			Type:       domain.LeaveTypeSick,
			StartDate:  time.Now().AddDate(0, 0, -7), // 7 days ago
			EndDate:    time.Now().AddDate(0, 0, -1), // 1 day ago
			Reason:     "Sick leave",
			Status:     domain.ReviewStatusRejected,
			Reviews: []domain.LeaveReview{
				{
					ReviewerID: employees[0].ID,
					Status:     domain.ReviewStatusRejected,
					Comment:    "Not enough evidence",
					ReviewedAt: getPtr(time.Now().AddDate(0, 0, -10)),
				},
			},
		},
		{
			EmployeeID:        employees[4].ID,
			Type:              domain.LeaveTypeAnnual,
			StartDate:         time.Now().AddDate(0, 0, 2),  // 2 days later
			EndDate:           time.Now().AddDate(0, 0, 20), // 20 days later
			Reason:            "Medical appointment",
			Status:            domain.ReviewStatusReviewing,
			CurrentReviewerID: getPtr(employees[2].ID),
			Reviews: []domain.LeaveReview{
				{
					ReviewerID: employees[2].ID,
					Status:     domain.ReviewStatusReviewing,
				},
			},
		},
	}

	for i := range leaves {
		leave := leaves[i]
		if err := r.db.WithContext(ctx).Create(&leave).Error; err != nil {
			return fmt.Errorf("failed to seed leave data: %w", err)
		}
	}

	return nil
}

func getPtr[T any](v T) *T {
	return &v
}
