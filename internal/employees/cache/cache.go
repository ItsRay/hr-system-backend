package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"hr-system/internal/cache"
	"hr-system/internal/common"
	"hr-system/internal/employees/domain"
)

type EmployeeCache interface {
	GetEmployeeByID(ctx context.Context, id int) (domain.Employee, error)
	SetEmployeeToCache(ctx context.Context, employee *domain.Employee, expiration time.Duration) error
	DeleteEmployeeCache(ctx context.Context, id int) error
	GetEmployees(ctx context.Context, page, pageSize int) (employees []domain.Employee, totalCount int, er error)
	SetEmployeesToCache(ctx context.Context, page, pageSize int, employees []domain.Employee, totalCount int, expiration time.Duration) error
	DeleteEmployeesListCache(ctx context.Context) error
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
	if errors.Is(err, redis.Nil) {
		return domain.Employee{}, common.ErrResourceNotFound
	} else if err != nil {
		return domain.Employee{}, err
	}

	if data == "" {
		return domain.Employee{}, common.ErrResourceNotFound
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

func (e *employeeCache) genEmployeesListCachePrefix() string {
	return fmt.Sprintf("%s_list", e.prefix)
}

func (e *employeeCache) genEmployeesListCacheKey(page, pageSize int) string {
	prefix := e.genEmployeesListCachePrefix()
	return fmt.Sprintf("%s_page_%d_page_size_%d", prefix, page, pageSize)
}

type EmployeesCacheData struct {
	Employees  []domain.Employee `json:"employees"`
	TotalCount int               `json:"total_count"`
}

func (e *employeeCache) GetEmployees(ctx context.Context, page, pageSize int) ([]domain.Employee, int, error) {
	cacheKey := e.genEmployeesListCacheKey(page, pageSize)
	data, err := e.cache.Get(ctx, cacheKey)
	if errors.Is(err, redis.Nil) {
		return nil, 0, common.ErrResourceNotFound
	} else if err != nil {
		return nil, 0, err
	}

	if data == "" {
		return nil, 0, common.ErrResourceNotFound
	}

	var employees EmployeesCacheData
	err = json.Unmarshal([]byte(data), &employees)
	if err != nil {
		return nil, 0, err
	}

	return employees.Employees, employees.TotalCount, nil
}

func (e *employeeCache) SetEmployeesToCache(ctx context.Context, page, pageSize int, employees []domain.Employee,
	totalCount int, expiration time.Duration) error {

	cacheData := EmployeesCacheData{
		Employees:  employees,
		TotalCount: totalCount,
	}

	cacheKey := e.genEmployeesListCacheKey(page, pageSize)

	jsonData, err := json.Marshal(cacheData)
	if err != nil {
		return err
	}

	return e.cache.Set(ctx, cacheKey, string(jsonData), expiration)
}

func (e *employeeCache) DeleteEmployeesListCache(ctx context.Context) error {
	return e.cache.DelByPrefix(ctx, e.genEmployeesListCachePrefix())
}
