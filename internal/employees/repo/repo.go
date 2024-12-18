package repo

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"hr-system/internal/employees/domain"
)

type EmployeeRepo interface {
	Create(employee *domain.Employee) error
	GetEmployeeByID(id int) (domain.Employee, error)
	GetEmployees() ([]domain.Employee, error)
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

	return domain.Employee{
		ID:          e.ID,
		Name:        e.Name,
		Email:       e.Email,
		Address:     e.Address,
		PhoneNumber: e.PhoneNumber,
		Positions:   domainPositions,
	}
}

func preloadPositions(db *gorm.DB) *gorm.DB {
	return db.Preload("Positions", func(db *gorm.DB) *gorm.DB {
		return db.Order("start_date DESC")
	})
}

func (r *employeeRepo) GetEmployeeByID(id int) (domain.Employee, error) {
	var employee Employee

	if err := preloadPositions(r.db).First(&employee, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Employee{}, domain.ErrResourceNotFound
		}
		return domain.Employee{}, err
	}

	return toDomainEmployee(&employee), nil
}

func (r *employeeRepo) GetEmployees() ([]domain.Employee, error) {
	var employeeModels []Employee

	err := preloadPositions(r.db).Find(&employeeModels).Error
	if err != nil {
		return nil, err
	}

	var employees []domain.Employee
	for i := range employeeModels {
		empModel := employeeModels[i]
		employee := domain.Employee{
			ID:          empModel.ID,
			Name:        empModel.Name,
			Email:       empModel.Email,
			Address:     empModel.Address,
			PhoneNumber: empModel.PhoneNumber,
		}

		var positions []domain.Position
		for j := range empModel.Positions {
			posModel := empModel.Positions[j]
			positions = append(positions, domain.Position{
				Title:        posModel.Title,
				Level:        posModel.Level,
				ManagerLevel: posModel.ManagerLevel,
				MonthSalary:  posModel.MonthSalary,
				StartDate:    posModel.StartDate,
				EndDate:      posModel.EndDate,
			})
		}

		employee.Positions = positions

		employees = append(employees, employee)
	}

	return employees, nil
}
