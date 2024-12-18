package repo

import (
	"context"
	"fmt"
	"time"
)

// SeedEmployees Seed function to seed Employee data
func (r *employeeRepo) SeedEmployees(ctx context.Context) error {
	// Create sample employees
	// 4 -> 2 -> 1 -> 0
	// 3 -> 0
	employees := []Employee{
		{
			Name:        "Alice Johnson",
			Email:       "alice.johnson@example.com",
			Address:     "123 Main St",
			PhoneNumber: "555-1234",
			ManagerID:   nil,
			Positions: []Position{
				{
					Title:        "Software Engineer",
					Level:        "Junior",
					ManagerLevel: 5,
					MonthSalary:  5000.00,
					StartDate:    time.Now().AddDate(0, -10, 0),
				},
			},
		},
		{
			Name:        "Bob Smith",
			Email:       "bob.smith@example.com",
			Address:     "456 Oak Rd",
			PhoneNumber: "555-5678",
			ManagerID:   nil,
			Positions: []Position{
				{
					Title:        "Senior Developer",
					Level:        "Senior",
					ManagerLevel: 4,
					MonthSalary:  8000.00,
					StartDate:    time.Now().AddDate(0, -6, 0),
				},
			},
		},
		{
			Name:        "Charlie Davis",
			Email:       "charlie.davis@example.com",
			Address:     "789 Pine St",
			PhoneNumber: "555-9101",
			ManagerID:   nil,
			Positions: []Position{
				{
					Title:        "HR Manager",
					Level:        "Manager",
					ManagerLevel: 3,
					MonthSalary:  7000.00,
					StartDate:    time.Now().AddDate(0, -4, 0),
				},
			},
		},
		{
			Name:        "David Lee",
			Email:       "david.lee@example.com",
			Address:     "101 Maple St",
			PhoneNumber: "555-1122",
			ManagerID:   nil,
			Positions: []Position{
				{
					Title:        "Product Manager",
					Level:        "Manager",
					ManagerLevel: 3,
					MonthSalary:  9500.00,
					StartDate:    time.Now().AddDate(0, -3, 0),
				},
			},
		},
		{
			Name:        "Eva Zhang",
			Email:       "eva.zhang@example.com",
			Address:     "202 Birch Rd",
			PhoneNumber: "555-3344",
			ManagerID:   nil,
			Positions: []Position{
				{
					Title:        "Data Scientist",
					Level:        "Junior",
					ManagerLevel: 0,
					MonthSalary:  6000.00,
					StartDate:    time.Now().AddDate(0, -1, 0), // 1 month ago
				},
			},
		},
	}

	// 1. Create Alice Johnson
	employee0 := employees[0]
	if err := r.db.WithContext(ctx).Create(&employee0).Error; err != nil {
		return fmt.Errorf("failed to seed Alice Johnson: %w", err)
	}

	// 2. Create Bob Smith and set Alice as its manager
	employee1 := employees[1]
	employee1.ManagerID = &employee0.ID
	if err := r.db.WithContext(ctx).Create(&employee1).Error; err != nil {
		return fmt.Errorf("failed to seed Bob Smith: %w", err)
	}

	// 3. Create Charlie Davis and set Bob Smith as its manager
	employee2 := employees[2]
	employee2.ManagerID = &employee1.ID
	if err := r.db.WithContext(ctx).Create(&employee2).Error; err != nil {
		return fmt.Errorf("failed to seed Charlie Davis: %w", err)
	}

	// 4. Create David Lee and set Alice as its manager
	employee3 := employees[3]
	employee3.ManagerID = &employee0.ID
	if err := r.db.WithContext(ctx).Create(&employee3).Error; err != nil {
		return fmt.Errorf("failed to seed David Lee: %w", err)
	}

	// 5. Create Eva Zhang and set Charlie Davis as its manager
	employee4 := employees[4]
	employee4.ManagerID = &employee2.ID
	if err := r.db.WithContext(ctx).Create(&employee4).Error; err != nil {
		return fmt.Errorf("failed to seed Eva Zhang: %w", err)
	}

	return nil
}
