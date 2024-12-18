package domain

import "time"

type Employee struct {
	ID          int        `json:"id"`
	Name        string     `json:"name" validate:"required"`
	Email       string     `json:"email" validate:"required"`
	Address     string     `json:"address"`
	PhoneNumber string     `json:"phone_number"`
	Positions   []Position `json:"positions" validate:"required,gt=0,dive"`
}

type Position struct {
	Title        string    `json:"title" validate:"required"`
	Level        string    `json:"level"`
	ManagerLevel int       `json:"manager_level"`
	MonthSalary  float64   `json:"month_salary" validate:"gte=0"`
	StartDate    time.Time `json:"start_date" validate:"required"`
	EndDate      time.Time `json:"end_time"`
}
