package repo

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"hr-system/internal/common"
	employee_repo "hr-system/internal/employees/repo"
	"hr-system/internal/leaves/domain"
)

type LeaveRepo interface {
	CreateLeave(ctx context.Context, leave *domain.Leave) error
	GetLeaveByID(ctx context.Context, id int) (domain.Leave, error)
	GetLeaves(ctx context.Context, query domain.LeaveQuery) ([]domain.Leave, error)
	UpdateLeaveAndReviews(ctx context.Context, leave *domain.Leave, reviews []domain.LeaveReview) error
}

type leaveRepo struct {
	db *gorm.DB
}

func NewLeaveRepo(ctx context.Context, db *gorm.DB, employeeRepo employee_repo.EmployeeRepo) (LeaveRepo, error) {
	repo := &leaveRepo{
		db: db,
	}

	if err := repo.ensureSchemaAndSeedData(ctx, employeeRepo); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *leaveRepo) ensureSchemaAndSeedData(ctx context.Context, employeeRepo employee_repo.EmployeeRepo) error {
	// migrate
	if err := r.db.AutoMigrate(domain.Leave{}); err != nil {
		return err
	}
	if err := r.db.AutoMigrate(domain.LeaveReview{}); err != nil {
		return err
	}

	// seed
	if err := r.SeedLeaveData(ctx, employeeRepo); err != nil {
		return err
	}
	return nil
}

func (r *leaveRepo) CreateLeave(ctx context.Context, leave *domain.Leave) error {
	if err := r.db.WithContext(ctx).Create(leave).Error; err != nil {
		return fmt.Errorf("failed to create leave: %w", err)
	}
	return nil
}

func preloadReviews(db *gorm.DB) *gorm.DB {
	return db.Preload("Reviews", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC")
	})
}

func (r *leaveRepo) GetLeaveByID(ctx context.Context, id int) (domain.Leave, error) {
	var leave domain.Leave
	db := r.db.WithContext(ctx)
	db = preloadReviews(db)
	if err := db.First(&leave, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Leave{}, common.ErrResourceNotFound
		}
		return domain.Leave{}, fmt.Errorf("failed to find leave with id %d: %w", id, err)
	}
	return leave, nil
}

func doTrans(db *gorm.DB, op func(*gorm.DB) error) (err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
		if err != nil {
			tx.Rollback()
		}
	}()

	if err := op(tx); err != nil {
		return err
	}
	return tx.Commit().Error
}

func (r *leaveRepo) UpdateLeaveAndReviews(ctx context.Context, leave *domain.Leave, reviews []domain.LeaveReview) error {
	err := doTrans(r.db.WithContext(ctx), func(tx *gorm.DB) error {
		result := tx.Save(leave)
		if result.Error != nil {
			return fmt.Errorf("failed to update leave: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("no rows affected, leave ID may not exist")
		}

		for _, review := range reviews {
			if review.ID == 0 {
				// if review.ID == 0, it's a new review
				result := tx.Create(&review)
				if result.Error != nil {
					return fmt.Errorf("failed to create leave review: %w", result.Error)
				}
			} else {
				// update existing review
				result := tx.Save(&review)
				if result.Error != nil {
					return fmt.Errorf("failed to update leave review: %w", result.Error)
				}
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update leave and reviews: %w", err)
	}

	return nil
}

func (r *leaveRepo) GetLeaves(ctx context.Context, query domain.LeaveQuery) ([]domain.Leave, error) {
	// TODO: pagination
	var leaves []domain.Leave

	db := r.db.WithContext(ctx).Preload("Reviews")
	if query.EmployeeID != nil {
		db = db.Where("employee_id = ?", *query.EmployeeID)
	}
	if query.CurrentReviewerID != nil {
		db = db.Where("current_reviewer_id = ?", *query.CurrentReviewerID)
	}
	if err := db.Order("id desc").Find(&leaves).Error; err != nil {
		return nil, fmt.Errorf("failed to get leaves: %w", err)
	}

	return leaves, nil
}
