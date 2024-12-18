package domain

import "time"

type LeaveType string

var (
	LeaveTypeAnnual LeaveType = "annual"
	LeaveTypeSick   LeaveType = "sick"
)

type ReviewStatus string

var (
	ReviewStatusReviewing ReviewStatus = "reviewing"
	ReviewStatusApproved  ReviewStatus = "approved"
	ReviewStatusRejected  ReviewStatus = "rejected"
)

type Leave struct {
	ID                int           `gorm:"primaryKey;autoIncrement"`
	EmployeeID        int           `gorm:"index:idx_employee_id" validate:"required"`
	Type              LeaveType     `gorm:"type:varchar(50);not null" validate:"required,oneof=annual sick"`
	StartDate         time.Time     `gorm:"type:date;not null" validate:"required"`
	EndDate           time.Time     `gorm:"type:date;not null" validate:"required"`
	Reason            string        `gorm:"type:varchar(255)"`
	Status            ReviewStatus  `gorm:"type:varchar(50);not null"`
	CurrentReviewerID *int          `gorm:"index:idx_current_reviewer_id"`
	Reviews           []LeaveReview `gorm:"foreignKey:LeaveID"`
	CreatedAt         time.Time     `gorm:"autoCreateTime"`
	UpdatedAt         time.Time     `gorm:"autoUpdateTime"`
}

type LeaveReview struct {
	ID         int          `gorm:"primaryKey;autoIncrement"`
	LeaveID    int          `gorm:"index:idx_leave_id"`
	ReviewerID int          `gorm:"type:int;not null"`
	Status     ReviewStatus `gorm:"type:varchar(50);not null"`
	Comment    string       `gorm:"type:varchar(255)"`
	ReviewedAt *time.Time   `gorm:"type:date"`
	CreatedAt  time.Time    `gorm:"autoCreateTime"`
	UpdatedAt  time.Time    `gorm:"autoUpdateTime"`
}

type LeaveQuery struct {
	EmployeeID        *int
	CurrentReviewerID *int
}
