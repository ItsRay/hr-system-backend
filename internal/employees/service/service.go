package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"hr-system/internal/employees/domain"
	"hr-system/internal/employees/repo"
)

type EmployeeService interface {
	CreateEmployee(employee *domain.Employee) (domain.Employee, error)
}

type employeeService struct {
	repo     repo.EmployeeRepo
	validate *validator.Validate
}

func NewEmployeeService(repo repo.EmployeeRepo) EmployeeService {
	return &employeeService{
		repo:     repo,
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

func (s *employeeService) CreateEmployee(employee *domain.Employee) (domain.Employee, error) {
	if err := s.validateCreateEmployee(employee); err != nil {
		return domain.Employee{}, err
	}

	err := s.repo.Create(employee)
	return *employee, err
}
