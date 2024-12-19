package repo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"hr-system/internal/employees/domain"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err != nil {
			t.Fatalf("failed to get database: %v", err)
		}
		sqlDB.Close()
	})
	return db
}

func TestEmployeeRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	repo, err := NewEmployeeRepo(db)
	assert.NoError(t, err)

	employee := &domain.Employee{
		Name:        "John Doe",
		Email:       "john.doe@example.com",
		Address:     "123 Main St",
		PhoneNumber: "123-456-7890",
		Positions: []domain.Position{
			{
				Title:       "Software Engineer",
				Level:       "Senior",
				MonthSalary: 10000,
				StartDate:   time.Now(),
			},
		},
	}

	err = repo.Create(context.Background(), employee)
	assert.NoError(t, err)
	assert.NotZero(t, employee.ID)
}

func TestEmployeeRepo_GetEmployeeByID(t *testing.T) {
	db := setupTestDB(t)
	repo, err := NewEmployeeRepo(db)
	assert.NoError(t, err)

	employee := &domain.Employee{
		Name:        "John Doe",
		Email:       "john.doe@example.com",
		Address:     "123 Main St",
		PhoneNumber: "123-456-7890",
		Positions: []domain.Position{
			{
				Title:       "Software Engineer",
				Level:       "Senior",
				MonthSalary: 10000,
				StartDate:   time.Now(),
			},
		},
	}

	err = repo.Create(context.Background(), employee)
	assert.NoError(t, err)

	fetchedEmployee, err := repo.GetEmployeeByID(context.Background(), employee.ID)
	assert.NoError(t, err)
	assert.Equal(t, employee.Name, fetchedEmployee.Name)
	assert.Equal(t, employee.Email, fetchedEmployee.Email)
}

func TestEmployeeRepo_GetEmployees(t *testing.T) {
	db := setupTestDB(t)
	repo, err := NewEmployeeRepo(db)
	assert.NoError(t, err)

	employee1 := &domain.Employee{
		Name:        "John Doe",
		Email:       "john.doe@example.com",
		Address:     "123 Main St",
		PhoneNumber: "123-456-7890",
		Positions: []domain.Position{
			{
				Title:       "Software Engineer",
				Level:       "Senior",
				MonthSalary: 10000,
				StartDate:   time.Now(),
			},
		},
	}

	employee2 := &domain.Employee{
		Name:        "Jane Smith",
		Email:       "jane.smith@example.com",
		Address:     "456 Elm St",
		PhoneNumber: "987-654-3210",
		Positions: []domain.Position{
			{
				Title:       "Product Manager",
				Level:       "Mid",
				MonthSalary: 8000,
				StartDate:   time.Now(),
			},
		},
	}

	err = repo.Create(context.Background(), employee1)
	assert.NoError(t, err)
	err = repo.Create(context.Background(), employee2)
	assert.NoError(t, err)

	employees, totalCount, err := repo.GetEmployees(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 2, totalCount)
	assert.Len(t, employees, 2)
}
