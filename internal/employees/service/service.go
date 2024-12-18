package service

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"

	"hr-system/internal/employees/cache"
	"hr-system/internal/employees/domain"
	"hr-system/internal/employees/repo"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, employee *domain.Employee) (domain.Employee, error)
	GetEmployeeByID(ctx context.Context, id int) (domain.Employee, error)
	GetEmployees(ctx context.Context, page, pageSize int) (employees []domain.Employee, totalCount int64, err error)
}

type employeeService struct {
	repo     repo.EmployeeRepo
	validate *validator.Validate
	cache    cache.EmployeeCache
}

func NewEmployeeService(repo repo.EmployeeRepo, cache cache.EmployeeCache) EmployeeService {
	return &employeeService{
		repo:     repo,
		cache:    cache,
		validate: validator.New(),
	}
}

func (s *employeeService) validateCreateEmployee(e *domain.Employee) error {
	if e == nil {
		return fmt.Errorf("employee is nil")
	}

	if err := s.validate.Struct(e); err != nil {
		return err
	}

	return nil
}

func (s *employeeService) CreateEmployee(ctx context.Context, employee *domain.Employee) (domain.Employee, error) {
	if err := s.validateCreateEmployee(employee); err != nil {
		return domain.Employee{}, err
	}

	err := s.repo.Create(ctx, employee)
	return *employee, err
}

func (s *employeeService) GetEmployeeByID(ctx context.Context, id int) (domain.Employee, error) {
	employee, err := s.repo.GetEmployeeByID(ctx, id)
	if err != nil {
		return domain.Employee{}, err
	}
	return employee, nil
}

func (s *employeeService) GetEmployees(ctx context.Context, page, pageSize int) ([]domain.Employee, int64, error) {
	employees, totalCount, err := s.repo.GetEmployees(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return employees, totalCount, nil
}
