package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"hr-system/internal/common"
	common_errors "hr-system/internal/common/errors"
	"hr-system/internal/employees/domain"
	mock_service "hr-system/internal/employees/service/mocks"
)

func genFakeEmployee() domain.Employee {
	return domain.Employee{
		ID:          1,
		Name:        "John Doe",
		Email:       "john.doe@example.com",
		Address:     "123 Main St",
		PhoneNumber: "1234567890",
		Positions: []domain.Position{
			{
				Title:        "Developer",
				Level:        "Senior",
				ManagerLevel: 2,
				MonthSalary:  5000.0,
				StartDate:    time.Now().Truncate(time.Second),
			},
		},
		ManagerID: nil,
	}
}

func TestCreateEmployee(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mock_service.NewEmployeeService(t)
	logger := common.NewLogger()
	handler := NewEmployeeHandler(logger, mockService)

	router := gin.Default()
	router.POST("/employees", handler.CreateEmployee)

	t.Run("success", func(t *testing.T) {
		employee := genFakeEmployee()
		reqBody := &CreateEmployeeRequest{
			Name:        employee.Name,
			Email:       employee.Email,
			Address:     employee.Address,
			PhoneNumber: employee.PhoneNumber,
			Position: Position{
				Title:        employee.Positions[0].Title,
				Level:        employee.Positions[0].Level,
				ManagerLevel: employee.Positions[0].ManagerLevel,
				MonthSalary:  employee.Positions[0].MonthSalary,
				StartDate:    employee.Positions[0].StartDate,
			},
			ManagerID: employee.ManagerID,
		}
		reqBodyBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/employees", bytes.NewBuffer(reqBodyBytes))

		employee.ID = 0
		mockService.On("CreateEmployee", mock.Anything, &employee).Return(employee, nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/employees", bytes.NewBuffer([]byte(`{}`)))

		mockService.On("CreateEmployee", mock.Anything, mock.Anything).Return(domain.Employee{}, common_errors.ErrInvalidInput)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetEmployeeByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mock_service.NewEmployeeService(t)
	logger := common.NewLogger()
	handler := NewEmployeeHandler(logger, mockService)

	router := gin.Default()
	router.GET("/employees/:id", handler.GetEmployeeByID)

	t.Run("success", func(t *testing.T) {
		employee := genFakeEmployee()

		mockService.On("GetEmployeeByID", mock.Anything, 1).Return(employee, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/employees/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid employee ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/employees/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("employee not found", func(t *testing.T) {
		mockService.On("GetEmployeeByID", mock.Anything, 1).
			Return(domain.Employee{}, common_errors.ErrResourceNotFound).Once()

		req, _ := http.NewRequest(http.MethodGet, "/employees/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetEmployees(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := mock_service.NewEmployeeService(t)
	logger := common.NewLogger()
	handler := NewEmployeeHandler(logger, mockService)

	router := gin.Default()
	router.GET("/employees", handler.GetEmployees)

	t.Run("success", func(t *testing.T) {
		employees := []domain.Employee{
			genFakeEmployee(),
		}

		mockService.On("GetEmployees", mock.Anything, 1, 10).Return(employees, 1, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/employees?page=1&page_size=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "1", w.Header().Get("X-Total-Count"))
		assert.Equal(t, "1", w.Header().Get("X-Page"))
		assert.Equal(t, "10", w.Header().Get("X-Page-Size"))
	})

	t.Run("invalid input", func(t *testing.T) {
		mockService.On("GetEmployees", mock.Anything, 1, 10).
			Return(nil, 0, common_errors.ErrInvalidInput).
			Once()

		req, _ := http.NewRequest(http.MethodGet, "/employees?page=1&page_size=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
