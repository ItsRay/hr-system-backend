package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"hr-system/internal/common"
	cache_mocks "hr-system/internal/employees/cache/mocks"
	"hr-system/internal/employees/domain"
	repo_mocks "hr-system/internal/employees/repo/mocks"
)

func newMockRepoAndCache(t *testing.T) (*repo_mocks.EmployeeRepo, *cache_mocks.EmployeeCache) {
	mockRepo := repo_mocks.NewEmployeeRepo(t)
	mockCache := cache_mocks.NewEmployeeCache(t)

	return mockRepo, mockCache
}

func genFakeEmployee() domain.Employee {
	startDate := time.Now()
	managerID := 1
	return domain.Employee{
		Name:        "John Doe",
		Email:       "john.doe@example.com",
		Address:     "123 Elm St",
		PhoneNumber: "123-456-7890",
		Positions: []domain.Position{
			{
				Title:        "Software Engineer",
				Level:        "Junior",
				ManagerLevel: 2,
				MonthSalary:  3500.00,
				StartDate:    startDate,
				EndDate:      nil, // No end date yet
			},
		},
		ManagerID: &managerID,
	}
}

func TestCreateEmployee(t *testing.T) {
	mockRepo, mockCache := newMockRepoAndCache(t)
	logger := common.NewLogger()
	service := NewEmployeeService(logger, mockRepo, mockCache)

	employee := genFakeEmployee()
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Employee")).Run(func(args mock.Arguments) {
		emp := args.Get(1).(*domain.Employee)
		emp.ID = 1 // Set the ID to match the expected employee
	}).Return(nil)
	mockCache.On("DeleteEmployeesListCache", mock.Anything).Return(nil)
	mockCache.On("SetEmployeeToCache", mock.Anything, &employee, 1*time.Hour).Return(nil)

	result, err := service.CreateEmployee(context.Background(), &employee)
	assert.NoError(t, err)
	assert.Equal(t, employee, result)
}

func TestGetEmployeeByID(t *testing.T) {
	mockRepo, mockCache := newMockRepoAndCache(t)
	logger := common.NewLogger()
	service := NewEmployeeService(logger, mockRepo, mockCache)

	employee := genFakeEmployee()
	employee.ID = 1
	mockCache.On("GetEmployeeByID", mock.Anything, employee.ID).Return(employee, nil)

	result, err := service.GetEmployeeByID(context.Background(), employee.ID)
	assert.NoError(t, err)
	assert.Equal(t, employee, result)
}

func TestGetEmployees(t *testing.T) {
	mockRepo, mockCache := newMockRepoAndCache(t)
	logger := common.NewLogger()
	service := NewEmployeeService(logger, mockRepo, mockCache)

	employees := []domain.Employee{
		genFakeEmployee(),
	}
	employees[0].ID = 1
	totalCount := 1

	mockCache.On("GetEmployees", mock.Anything, 1, 10).Return(employees, totalCount, nil)

	result, count, err := service.GetEmployees(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, employees, result)
	assert.Equal(t, totalCount, count)
}
