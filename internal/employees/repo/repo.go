package repo

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"hr-system/internal/employees/domain"
)

type EmployeeRepo interface {
	Create(employee *domain.Employee) error
	GetByID(id int) (*domain.Employee, error)
}

type Employee struct {
	ID          int        `gorm:"primaryKey;autoIncrement"`
	Name        string     `gorm:"type:varchar(255);not null"`
	Email       string     `gorm:"type:varchar(255);unique;not null"`
	Address     string     `gorm:"type:varchar(255)"`
	PhoneNumber string     `gorm:"type:varchar(20)"`
	Positions   []Position `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Position struct {
	ID           int       `gorm:"primaryKey;autoIncrement"`
	EmployeeID   int       `gorm:"index"`
	Title        string    `gorm:"type:varchar(255)"`
	Level        string    `gorm:"type:varchar(50)"`
	ManagerLevel int       `gorm:"type:int;default:0"`
	MonthSalary  float64   `gorm:"type:decimal(10,2);not null"`
	StartDate    time.Time `gorm:"type:date;not null"`
	EndDate      time.Time `gorm:"type:date"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

type employeeRepo struct {
	db *gorm.DB
}

func ToRepositoryEmployee(e *domain.Employee) Employee {
	return Employee{
		Name:        e.Name,
		Email:       e.Email,
		Address:     e.Address,
		PhoneNumber: e.PhoneNumber,
	}
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
	return r.db.AutoMigrate(Employee{})
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

func (r *employeeRepo) Create(e *domain.Employee) error {
	tx := r.db.Begin()

	employee := &Employee{
		Name:        e.Name,
		Email:       e.Email,
		Address:     e.Address,
		PhoneNumber: e.PhoneNumber,
		// TODO: Check if this works
		Positions: []Position{toRepoPosition(0, &e.Positions[0])},
	}
	if err := tx.Create(employee).Error; err != nil {
		tx.Rollback()
		return err
	}

	// TODO: positions
	position := toRepoPosition(0, &e.Positions[0])
	if err := tx.Create(&position).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *employeeRepo) GetByID(id int) (*domain.Employee, error) {
	var repoEmployee Employee

	if err := r.db.First(&repoEmployee, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrEmployeeNotFound
		}
		return nil, err
	}

	// 將資料庫層的 Employee 映射到 domain 層的 Employee 並返回
	return &domain.Employee{
		ID:          repoEmployee.ID,
		Name:        repoEmployee.Name,
		Email:       repoEmployee.Email,
		Address:     repoEmployee.Address,
		PhoneNumber: repoEmployee.PhoneNumber,
	}, nil
}
