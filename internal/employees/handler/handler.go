package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"hr-system/internal/common"
	"hr-system/internal/employees/domain"
	"hr-system/internal/employees/service"
	"hr-system/internal/middleware"
)

type EmployeeHandler struct {
	service  service.EmployeeService
	validate *validator.Validate
	logger   *common.Logger
}

func NewEmployeeHandler(logger *common.Logger, service service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		service:  service,
		validate: validator.New(),
		logger:   logger,
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
	ctx := c.Request.Context()

	var req *CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, middleware.CreateErrResp("invalid request body, cause: %s", err))
		return
	}
	h.logger.Info("create employee request: %+v", req)

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, middleware.CreateErrResp("invalid request body, cause: %s", err))
		return
	}
	h.logger.Info("CreateEmployee pass validate: %+v", req)

	employee, err := h.service.CreateEmployee(ctx, &domain.Employee{
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
		c.JSON(http.StatusInternalServerError, middleware.CreateErrResp("failed to create employee, cause: %s\n", err))
		return
	}

	c.JSON(http.StatusCreated, &employee)
}

func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	employee, err := h.service.GetEmployeeByID(ctx, id)
	if err != nil {
		if errors.Is(err, common.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) GetEmployees(c *gin.Context) {
	// TODO: order by query param
	ctx := c.Request.Context()

	var page, pageSize int

	if p := c.DefaultQuery("page", "1"); p != "" {
		page, _ = strconv.Atoi(p)
	}
	if ps := c.DefaultQuery("page_size", "10"); ps != "" {
		pageSize, _ = strconv.Atoi(ps)
	}

	employees, totalCount, err := h.service.GetEmployees(ctx, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 設置響應頭
	c.Header("X-Total-Count", strconv.Itoa(int(totalCount)))
	c.Header("X-Page", strconv.Itoa(page))
	c.Header("X-Page-Size", strconv.Itoa(pageSize))

	c.JSON(http.StatusOK, employees)
}
