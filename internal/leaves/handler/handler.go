package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"hr-system/internal/common"
	"hr-system/internal/leaves/domain"
	"hr-system/internal/leaves/service"
	"hr-system/internal/middleware"
)

type LeaveHandler struct {
	leaveService service.LeaveService
	logger       *common.Logger
}

func NewLeaveHandler(logger *common.Logger, leaveService service.LeaveService) *LeaveHandler {
	return &LeaveHandler{
		leaveService: leaveService,
		logger:       logger,
	}
}

type CreateLeaveRequest struct {
	EmployeeID int       `json:"employee_id" binding:"required"`
	Type       string    `json:"type" binding:"required"`
	StartDate  time.Time `json:"start_date" binding:"required"`
	EndDate    time.Time `json:"end_date" binding:"required"`
	Reason     string    `json:"reason"`
}

func (h *LeaveHandler) CreateLeave(c *gin.Context) {
	// TODO: check employee rest leaves & leave date conflict
	ctx := c.Request.Context()

	var req CreateLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, middleware.CreateErrResp("invalid request body, cause: %v", err))
		return
	}

	leave, err := h.leaveService.CreateLeave(ctx, &domain.Leave{
		EmployeeID: req.EmployeeID,
		Type:       domain.LeaveType(req.Type),
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Reason:     req.Reason,
	})
	if err != nil {
		if errors.Is(err, common.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, middleware.CreateErrResp("employee not found"))
		} else if errors.Is(err, common.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, middleware.CreateErrResp("invalid input, cause: %v", err))
		} else {
			c.JSON(http.StatusInternalServerError, middleware.CreateErrResp("Failed to create leave: %v", err))
		}
		return
	}

	c.JSON(http.StatusCreated, &leave)
}

type ReviewLeaveRequest struct {
	ReviewerID int                 `json:"reviewer_id" binding:"required"`
	Decision   domain.ReviewStatus `json:"decision" binding:"required,oneof=approved rejected"`
	Comment    string              `json:"comment"`
}

func (h *LeaveHandler) ReviewLeave(c *gin.Context) {
	ctx := c.Request.Context()

	leaveID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, middleware.CreateErrResp("Invalid leave ID"))
		return
	}

	var req ReviewLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind review leave request: %v", err)
		c.JSON(http.StatusBadRequest, middleware.CreateErrResp("invalid request body, cause: %v", err))
		return
	}

	// TODO: Use JWT to get the reviewer ID
	err = h.leaveService.ReviewLeave(ctx, leaveID, req.ReviewerID, req.Decision, req.Comment)
	if err != nil {
		if errors.Is(err, common.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, middleware.CreateErrResp("leave not found, cause: %v", err))
		} else if errors.Is(err, common.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, middleware.CreateErrResp("invalid input, cause: %v", err))
		} else {
			c.JSON(http.StatusInternalServerError, middleware.CreateErrResp("failed to review leave, cause: %v", err))
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *LeaveHandler) GetLeaves(c *gin.Context) {
	ctx := c.Request.Context()

	query := domain.LeavesQuery{}
	employeeID := c.Query("employee_id")
	if employeeID != "" {
		id, err := strconv.Atoi(employeeID)
		if err != nil {
			c.JSON(http.StatusBadRequest, middleware.CreateErrResp("invalid employee_id"))
			return
		}
		query.EmployeeID = &id
	}

	reviewerID := c.Query("current_reviewer_id")
	if reviewerID != "" {
		id, err := strconv.Atoi(reviewerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, middleware.CreateErrResp("invalid current_reviewer_id"))
			return
		}
		query.CurrentReviewerID = &id
	}

	leaves, err := h.leaveService.GetLeaves(ctx, query)
	if err != nil {
		if errors.Is(err, common.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, middleware.CreateErrResp("invalid input, cause: %v", err))
		} else {
			c.JSON(http.StatusInternalServerError, middleware.CreateErrResp("failed to get leaves: %v", err))
		}
		return
	}

	c.JSON(http.StatusOK, leaves)
}

func (h *LeaveHandler) GetLeaveByID(c *gin.Context) {
	ctx := c.Request.Context()

	leaveID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, middleware.CreateErrResp("Invalid leave ID"))
		return
	}

	leave, err := h.leaveService.GetLeaveByID(ctx, leaveID)
	if err != nil {
		if errors.Is(err, common.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, middleware.CreateErrResp("leave not found, cause: %v", err))
		} else {
			c.JSON(http.StatusInternalServerError, middleware.CreateErrResp("failed to get leave: %v", err))
		}
		return
	}

	c.JSON(http.StatusOK, leave)
}
