package cache

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"hr-system/internal/cache"
	"hr-system/internal/employees/domain"
)

func setupTestRedis() (*miniredis.Miniredis, *cache.Cache) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return mr, cache.NewCache(rdb)
}

func TestEmployeeCache_GetEmployeeByID(t *testing.T) {
	mr, c := setupTestRedis()
	defer mr.Close()

	employeeCache := NewEmployeeCache(c, "test")

	employee := domain.Employee{ID: 1, Name: "John Doe"}
	data, _ := json.Marshal(employee)
	mr.Set("test_id_1", string(data))

	ctx := context.Background()
	result, err := employeeCache.GetEmployeeByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, employee, result)
}

func TestEmployeeCache_SetEmployeeToCache(t *testing.T) {
	mr, c := setupTestRedis()
	defer mr.Close()

	employeeCache := NewEmployeeCache(c, "test")

	employee := &domain.Employee{ID: 1, Name: "John Doe"}
	ctx := context.Background()
	err := employeeCache.SetEmployeeToCache(ctx, employee, time.Minute)

	assert.NoError(t, err)

	data, err := c.Get(ctx, "test_id_1")
	assert.NoError(t, err)

	var cachedEmployee domain.Employee
	err = json.Unmarshal([]byte(data), &cachedEmployee)
	assert.NoError(t, err)
	assert.Equal(t, employee, &cachedEmployee)
}

func TestEmployeeCache_DeleteEmployeeCache(t *testing.T) {
	mr, c := setupTestRedis()
	defer mr.Close()

	employeeCache := NewEmployeeCache(c, "test")

	employee := &domain.Employee{ID: 1, Name: "John Doe"}
	ctx := context.Background()
	err := employeeCache.SetEmployeeToCache(ctx, employee, time.Minute)
	assert.NoError(t, err)

	err = employeeCache.DeleteEmployeeCache(ctx, 1)
	assert.NoError(t, err)

	_, err = c.Get(ctx, "test_id_1")
	assert.Equal(t, redis.Nil, err)
}

func TestEmployeeCache_GetEmployees(t *testing.T) {
	mr, c := setupTestRedis()
	defer mr.Close()

	employeeCache := NewEmployeeCache(c, "test")

	employees := []domain.Employee{
		{ID: 1, Name: "John Doe"},
		{ID: 2, Name: "Jane Doe"},
	}
	cacheData := EmployeesCacheData{
		Employees:  employees,
		TotalCount: 2,
	}
	jsonData, _ := json.Marshal(cacheData)
	mr.Set("test_list_page_1_page_size_2", string(jsonData))

	ctx := context.Background()
	result, totalCount, err := employeeCache.GetEmployees(ctx, 1, 2)

	assert.NoError(t, err)
	assert.Equal(t, employees, result)
	assert.Equal(t, 2, totalCount)
}

func TestEmployeeCache_SetEmployeesToCache(t *testing.T) {
	_, c := setupTestRedis()

	employeeCache := NewEmployeeCache(c, "test")

	employees := []domain.Employee{
		{ID: 1, Name: "John Doe"},
		{ID: 2, Name: "Jane Doe"},
	}
	ctx := context.Background()
	err := employeeCache.SetEmployeesToCache(ctx, 1, 2, employees, 2, time.Minute)
	assert.NoError(t, err)

	data, err := c.Get(ctx, "test_list_page_1_page_size_2")
	assert.NoError(t, err)

	var cachedData EmployeesCacheData
	err = json.Unmarshal([]byte(data), &cachedData)
	assert.NoError(t, err)
	assert.Equal(t, employees, cachedData.Employees)
	assert.Equal(t, 2, cachedData.TotalCount)
}

func TestEmployeeCache_DeleteEmployeesListCache(t *testing.T) {
	mr, c := setupTestRedis()
	defer mr.Close()

	employeeCache := NewEmployeeCache(c, "test")

	employees := []domain.Employee{
		{ID: 1, Name: "John Doe"},
		{ID: 2, Name: "Jane Doe"},
	}
	cacheData := EmployeesCacheData{
		Employees:  employees,
		TotalCount: 2,
	}
	jsonData, _ := json.Marshal(cacheData)
	mr.Set("test_list_page_1_page_size_2", string(jsonData))

	ctx := context.Background()
	err := employeeCache.DeleteEmployeesListCache(ctx)
	assert.NoError(t, err)

	_, err = c.Get(ctx, "test_list_page_1_page_size_2")
	assert.Equal(t, redis.Nil, err)
}
