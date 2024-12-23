// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "hr-system/internal/employees/domain"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// EmployeeCache is an autogenerated mock type for the EmployeeCache type
type EmployeeCache struct {
	mock.Mock
}

// DeleteEmployeeCache provides a mock function with given fields: ctx, id
func (_m *EmployeeCache) DeleteEmployeeCache(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteEmployeeCache")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteEmployeesListCache provides a mock function with given fields: ctx
func (_m *EmployeeCache) DeleteEmployeesListCache(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteEmployeesListCache")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetEmployeeByID provides a mock function with given fields: ctx, id
func (_m *EmployeeCache) GetEmployeeByID(ctx context.Context, id int) (domain.Employee, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetEmployeeByID")
	}

	var r0 domain.Employee
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (domain.Employee, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) domain.Employee); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.Employee)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEmployees provides a mock function with given fields: ctx, page, pageSize
func (_m *EmployeeCache) GetEmployees(ctx context.Context, page int, pageSize int) ([]domain.Employee, int, error) {
	ret := _m.Called(ctx, page, pageSize)

	if len(ret) == 0 {
		panic("no return value specified for GetEmployees")
	}

	var r0 []domain.Employee
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) ([]domain.Employee, int, error)); ok {
		return rf(ctx, page, pageSize)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int) []domain.Employee); ok {
		r0 = rf(ctx, page, pageSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Employee)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int) int); ok {
		r1 = rf(ctx, page, pageSize)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, int, int) error); ok {
		r2 = rf(ctx, page, pageSize)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SetEmployeeToCache provides a mock function with given fields: ctx, employee, expiration
func (_m *EmployeeCache) SetEmployeeToCache(ctx context.Context, employee *domain.Employee, expiration time.Duration) error {
	ret := _m.Called(ctx, employee, expiration)

	if len(ret) == 0 {
		panic("no return value specified for SetEmployeeToCache")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Employee, time.Duration) error); ok {
		r0 = rf(ctx, employee, expiration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetEmployeesToCache provides a mock function with given fields: ctx, page, pageSize, employees, totalCount, expiration
func (_m *EmployeeCache) SetEmployeesToCache(ctx context.Context, page int, pageSize int, employees []domain.Employee, totalCount int, expiration time.Duration) error {
	ret := _m.Called(ctx, page, pageSize, employees, totalCount, expiration)

	if len(ret) == 0 {
		panic("no return value specified for SetEmployeesToCache")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, []domain.Employee, int, time.Duration) error); ok {
		r0 = rf(ctx, page, pageSize, employees, totalCount, expiration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewEmployeeCache creates a new instance of EmployeeCache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEmployeeCache(t interface {
	mock.TestingT
	Cleanup(func())
}) *EmployeeCache {
	mock := &EmployeeCache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
