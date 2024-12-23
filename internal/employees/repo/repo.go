package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	common_errors "hr-system/internal/common/errors"
	"hr-system/internal/employees/domain"
)

type EmployeeRepo interface {
	// TODO: find a better place to put this
	SeedData(ctx context.Context) error
	Create(ctx context.Context, employee *domain.Employee) error
	GetEmployeeByID(ctx context.Context, id int) (domain.Employee, error)
	GetEmployees(ctx context.Context, page, pageSize int) (employees []domain.Employee, totalCount int, err error)
}

type Employee struct {
	ID          int        `gorm:"primaryKey;autoIncrement"`
	Name        string     `gorm:"type:varchar(255);not null"`
	Email       string     `gorm:"type:varchar(255);unique;not null"`
	Address     string     `gorm:"type:varchar(255)"`
	PhoneNumber string     `gorm:"type:varchar(20)"`
	ManagerID   *int       `gorm:"index:idx_manager_id"`
	Manager     *Employee  `gorm:"foreignKey:ManagerID;constraint:OnDelete:SET NULL"`
	Positions   []Position `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Position struct {
	ID           int        `gorm:"primaryKey;autoIncrement"`
	EmployeeID   int        `gorm:"index"`
	Title        string     `gorm:"type:varchar(255)"`
	Level        string     `gorm:"type:varchar(50)"`
	ManagerLevel int        `gorm:"type:int;default:0"`
	MonthSalary  float64    `gorm:"type:decimal(10,2);not null"`
	StartDate    time.Time  `gorm:"type:date;not null"`
	EndDate      *time.Time `gorm:"type:date"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
}

type employeeRepo struct {
	db *gorm.DB
}

func NewEmployeeRepo(db *gorm.DB) (EmployeeRepo, error) {
	repo := &employeeRepo{
		db: db,
	}

	if err := repo.ensureSchema(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *employeeRepo) ensureSchema() error {
	// AutoMigrate
	if err := r.db.AutoMigrate(Employee{}); err != nil {
		return err
	}
	if err := r.db.AutoMigrate(Position{}); err != nil {
		return err
	}

	return nil
}

func (r *employeeRepo) SeedData(ctx context.Context) error {
	return r.SeedEmployees(ctx)
}

func toRepoPosition(employeeID int, p *domain.Position) Position {
	return Position{
		EmployeeID:   employeeID,
		Title:        p.Title,
		Level:        p.Level,
		ManagerLevel: p.ManagerLevel,
		MonthSalary:  p.MonthSalary,
		StartDate:    p.StartDate,
	}
}

func (r *employeeRepo) Create(ctx context.Context, e *domain.Employee) error {
	positions := make([]Position, 0, len(e.Positions))
	for i := range e.Positions {
		positions = append(positions, toRepoPosition(0, &e.Positions[i]))
	}
	employee := &Employee{
		Name:        e.Name,
		Email:       e.Email,
		Address:     e.Address,
		PhoneNumber: e.PhoneNumber,
		ManagerID:   e.ManagerID,
		Positions:   positions,
	}
	if err := r.db.WithContext(ctx).Create(employee).Error; err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}
	e.ID = employee.ID

	return nil
}

func toDomainEmployee(e *Employee) domain.Employee {
	var domainPositions []domain.Position
	for i := range e.Positions {
		p := e.Positions[i]
		domainPositions = append(domainPositions, domain.Position{
			Title:        p.Title,
			Level:        p.Level,
			ManagerLevel: p.ManagerLevel,
			MonthSalary:  p.MonthSalary,
			StartDate:    p.StartDate,
			EndDate:      p.EndDate,
		})
	}

	var manager *domain.Employee
	if e.Manager != nil {
		managerPtr := toDomainEmployee(e.Manager)
		manager = &managerPtr
	}
	return domain.Employee{
		ID:          e.ID,
		Name:        e.Name,
		Email:       e.Email,
		Address:     e.Address,
		PhoneNumber: e.PhoneNumber,
		Positions:   domainPositions,
		ManagerID:   e.ManagerID,
		Manager:     manager,
	}
}

func preloadPositions(db *gorm.DB) *gorm.DB {
	return db.Preload("Positions", func(db *gorm.DB) *gorm.DB {
		return db.Order("start_date DESC")
	})
}

func preloadManager(db *gorm.DB) *gorm.DB {
	return db.Preload("Manager")
}

func (r *employeeRepo) GetEmployeeByID(ctx context.Context, id int) (domain.Employee, error) {
	var employee Employee

	db := r.db.WithContext(ctx)
	db = preloadPositions(db)
	db = preloadManager(db)
	if err := db.First(&employee, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Employee{}, common_errors.ErrResourceNotFound
		}
		return domain.Employee{}, err
	}

	return toDomainEmployee(&employee), nil
}

func (r *employeeRepo) GetEmployees(ctx context.Context, page, pageSize int) ([]domain.Employee, int, error) {
	var employeeModels []Employee

	offset := (page - 1) * pageSize

	var totalCountInt64 int64
	if err := r.db.WithContext(ctx).Model(&Employee{}).Count(&totalCountInt64).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employees: %w", err)
	}
	totalCount := int(totalCountInt64)

	db := r.db.WithContext(ctx)
	db = preloadPositions(db)
	err := db.Limit(pageSize).Offset(offset).Find(&employeeModels).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get employees: %w", err)
	}

	employees := make([]domain.Employee, 0, len(employeeModels))
	for i := range employeeModels {
		employees = append(employees, toDomainEmployee(&employeeModels[i]))
	}

	return employees, totalCount, nil
}
