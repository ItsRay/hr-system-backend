package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"

	"hr-system/internal/common"
	"hr-system/internal/employees/cache"
	"hr-system/internal/employees/domain"
	"hr-system/internal/employees/repo"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, employee *domain.Employee) (domain.Employee, error)
	GetEmployeeByID(ctx context.Context, id int) (domain.Employee, error)
	GetEmployees(ctx context.Context, page, pageSize int) (employees []domain.Employee, totalCount int, err error)
}

type employeeService struct {
	repo     repo.EmployeeRepo
	validate *validator.Validate
	cache    cache.EmployeeCache
	logger   *common.Logger
}

func NewEmployeeService(logger *common.Logger, repo repo.EmployeeRepo, cache cache.EmployeeCache) EmployeeService {
	return &employeeService{
		repo:     repo,
		cache:    cache,
		validate: validator.New(),
		logger:   logger,
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
		return domain.Employee{}, fmt.Errorf("%w, detail: %s", common.ErrInvalidInput, err)
	}

	// TODO: check manager level vs manager
	err := s.repo.Create(ctx, employee)
	if err != nil {
		return domain.Employee{}, err
	}

	err = s.cache.DeleteEmployeesListCache(ctx)
	if err != nil {
		s.logger.Errorf("failed to update cache, cause: %s", err)
	}

	err = s.cache.SetEmployeeToCache(ctx, employee, 1*time.Hour)
	if err != nil {
		s.logger.Warnf("failed to update cache, cause: %s", err)
	}

	return *employee, nil
}

func (s *employeeService) GetEmployeeByID(ctx context.Context, id int) (domain.Employee, error) {
	employee, err := s.cache.GetEmployeeByID(ctx, id)
	if err == nil {
		s.logger.Infof("[Cache Hit] employee id: %d", id)
		return employee, nil
	}
	if !errors.Is(err, common.ErrResourceNotFound) {
		s.logger.Warnf("failed to get employee from cache, cause: %s", err)
	}

	employee, err = s.repo.GetEmployeeByID(ctx, id)
	if err != nil {
		if errors.Is(err, common.ErrResourceNotFound) {
			return domain.Employee{}, common.ErrResourceNotFound
		}
		return domain.Employee{}, err
	}

	if err = s.cache.SetEmployeeToCache(ctx, &employee, 1*time.Hour); err != nil {
		s.logger.Warnf("failed to update cache, cause: %s", err)
	}

	return employee, nil
}

func (s *employeeService) GetEmployees(ctx context.Context, page, pageSize int) ([]domain.Employee, int, error) {
	if page < 1 || pageSize < 1 {
		return nil, 0, fmt.Errorf("%w, invalid page(%d) or page size(%d)", common.ErrInvalidInput, page, pageSize)
	}
	employees, totalCount, err := s.cache.GetEmployees(ctx, page, pageSize)
	if err == nil {
		s.logger.Infof("[Cache Hit] emplyees page: %d, pageSize: %d", page, pageSize)
		return employees, totalCount, nil
	}
	if err != nil && !errors.Is(err, common.ErrResourceNotFound) {
		s.logger.Warnf("failed to get employees from cache, cause: %s", err)
	}

	employees, totalCount, err = s.repo.GetEmployees(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	if err := s.cache.SetEmployeesToCache(ctx, page, pageSize, employees, totalCount, 1*time.Hour); err != nil {
		s.logger.Warnf("failed to update cache, cause: %s", err)
	}

	return employees, totalCount, nil
}
