package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"hr-system/internal/cache"
	"hr-system/internal/employees/domain"
)

type EmployeeCache interface {
	GetEmployeeByID(ctx context.Context, id int) (domain.Employee, error)
	SetEmployeeToCache(ctx context.Context, employee *domain.Employee, expiration time.Duration) error
	DeleteEmployeeCache(ctx context.Context, id int) error
	GetEmployees(ctx context.Context, page, pageSize int) ([]domain.Employee, error)
	SetEmployeesToCache(ctx context.Context, employees []domain.Employee, page, pageSize int, expiration time.Duration) error
	DeleteEmployeesCache(ctx context.Context, page, pageSize int) error
}

func NewEmployeeCache(c *cache.Cache, prefix string) EmployeeCache {
	return &employeeCache{
		cache:  c,
		prefix: prefix,
	}
}

type employeeCache struct {
	cache  *cache.Cache
	prefix string
}

func (e *employeeCache) genEmployeeCacheKey(prefix string, id int) string {
	return fmt.Sprintf("%s_id_%d", prefix, id)
}

func (e *employeeCache) GetEmployeeByID(ctx context.Context, id int) (domain.Employee, error) {
	cacheKey := e.genEmployeeCacheKey(e.prefix, id)
	data, err := e.cache.Get(ctx, cacheKey)
	if err != nil {
		return domain.Employee{}, err
	}

	if data == "" {
		return domain.Employee{}, nil
	}

	var employee domain.Employee
	err = json.Unmarshal([]byte(data), &employee)
	if err != nil {
		return domain.Employee{}, err
	}

	return employee, nil
}

func (e *employeeCache) SetEmployeeToCache(ctx context.Context, employee *domain.Employee, expiration time.Duration) error {
	if employee == nil {
		return errors.New("employee is nil")
	}
	cacheKey := e.genEmployeeCacheKey(e.prefix, employee.ID)
	jsonData, err := json.Marshal(employee)
	if err != nil {
		return fmt.Errorf("marshal employee error: %w", err)
	}
	return e.cache.Set(ctx, cacheKey, string(jsonData), expiration)
}

func (e *employeeCache) DeleteEmployeeCache(ctx context.Context, id int) error {
	cacheKey := e.genEmployeeCacheKey(e.prefix, id)
	return e.cache.Del(ctx, cacheKey)
}

func (e *employeeCache) genEmployeesListCacheKey(prefix string, page, pageSize int) string {
	return fmt.Sprintf("%s_list_page_%d_page_size_%d", prefix, page, pageSize)
}

func (e *employeeCache) GetEmployees(ctx context.Context, page, pageSize int) ([]domain.Employee, error) {
	cacheKey := e.genEmployeesListCacheKey(e.prefix, page, pageSize)
	data, err := e.cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	if data == "" {
		return nil, nil
	}

	var employees []domain.Employee
	err = json.Unmarshal([]byte(data), &employees)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (e *employeeCache) SetEmployeesToCache(ctx context.Context, employees []domain.Employee, page, pageSize int, expiration time.Duration) error {
	cacheKey := e.genEmployeesListCacheKey(e.prefix, page, pageSize)
	jsonData, err := json.Marshal(employees)
	if err != nil {
		return err
	}
	return e.cache.Set(ctx, cacheKey, string(jsonData), expiration)
}

func (e *employeeCache) DeleteEmployeesCache(ctx context.Context, page, pageSize int) error {
	cacheKey := e.genEmployeesListCacheKey(e.prefix, page, pageSize)
	return e.cache.Del(ctx, cacheKey)
}
