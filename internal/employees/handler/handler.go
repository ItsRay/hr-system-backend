package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"hr-system/internal/employees/domain"
	"hr-system/internal/employees/service"
)

type EmployeeHandler struct {
	service  service.EmployeeService
	validate *validator.Validate
}

func NewEmployeeHandler(service service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		service:  service,
		validate: validator.New(),
	}
}

type CreateEmployeeRequest struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Address     string   `json:"address"`
	PhoneNumber string   `json:"phone_number"`
	Position    Position `json:"position_level"`
}

type Position struct {
	Title        string    `json:"title"`
	Level        string    `json:"level"`
	ManagerLevel int       `json:"manager_level"`
	MonthSalary  float64   `json:"month_salary"`
	StartDate    time.Time `json:"start_date"`
}

type EmployeeResponse struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req *CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("invalid request body, cause: %w", err))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("invalid request body, cause: %w", err))
		return
	}

	employee, err := h.service.CreateEmployee(&domain.Employee{
		Name:        req.Name,
		Email:       req.Email,
		Address:     req.Address,
		PhoneNumber: req.PhoneNumber,
		Positions: []domain.Position{
			{
				Title:        req.Position.Title,
				Level:        req.Position.Level,
				ManagerLevel: req.Position.ManagerLevel,
				MonthSalary:  req.Position.MonthSalary,
				StartDate:    req.Position.StartDate,
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to create employee, cause: %w\n", err))
		return
	}

	c.JSON(http.StatusCreated, &employee)
}

func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	employee, err := h.service.GetEmployeeByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) GetEmployees(c *gin.Context) {
	employees, err := h.service.GetEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employees)
}
