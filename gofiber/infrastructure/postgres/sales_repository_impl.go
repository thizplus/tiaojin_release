package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type SalesRepositoryImpl struct {
	db *gorm.DB
}

func NewSalesRepository(db *gorm.DB) repositories.SalesRepository {
	return &SalesRepositoryImpl{db: db}
}

func (r *SalesRepositoryImpl) Create(ctx context.Context, entity *models.Sales) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create sales: %w", err)
	}
	return nil
}

func (r *SalesRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Sales, error) {
	var sales models.Sales
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&sales).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get sales by ID: %w", err)
	}

	return &sales, nil
}

func (r *SalesRepositoryImpl) Update(ctx context.Context, id uuid.UUID, entity *models.Sales) error {
	// Use Select("*") to update all fields including zero values like false
	result := r.db.WithContext(ctx).Where("id = ?", id).Select("*").Updates(entity)

	if result.Error != nil {
		return fmt.Errorf("failed to update sales: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("sales with ID %s not found", id)
	}

	return nil
}

func (r *SalesRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Sales{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete sales: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("sales with ID %s not found", id)
	}

	return nil
}

func (r *SalesRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Sales, error) {
	var sales []*models.Sales

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&sales).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list sales: %w", err)
	}

	return sales, nil
}

func (r *SalesRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Sales{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count sales: %w", err)
	}

	return count, nil
}

func (r *SalesRepositoryImpl) GetByCode(ctx context.Context, code string) (*models.Sales, error) {
	var sales models.Sales
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&sales).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get sales by code: %w", err)
	}

	return &sales, nil
}

func (r *SalesRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Sales, error) {
	var sales models.Sales
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		First(&sales).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get sales by user ID: %w", err)
	}

	return &sales, nil
}

func (r *SalesRepositoryImpl) GetActive(ctx context.Context, offset, limit int) ([]*models.Sales, error) {
	var sales []*models.Sales

	err := r.db.WithContext(ctx).
		Where("active = ?", true).
		Order("first_name ASC, last_name ASC").
		Offset(offset).
		Limit(limit).
		Find(&sales).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get active sales: %w", err)
	}

	return sales, nil
}

func (r *SalesRepositoryImpl) SearchByName(ctx context.Context, name string, offset, limit int) ([]*models.Sales, error) {
	var sales []*models.Sales
	searchPattern := "%" + name + "%"

	err := r.db.WithContext(ctx).
		Where("first_name ILIKE ? OR last_name ILIKE ? OR nick_name ILIKE ? OR code ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern).
		Order("first_name ASC, last_name ASC").
		Offset(offset).
		Limit(limit).
		Find(&sales).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search sales by name: %w", err)
	}

	return sales, nil
}

func (r *SalesRepositoryImpl) SearchByKeyword(ctx context.Context, keyword string, offset, limit int) ([]*models.Sales, error) {
	var sales []*models.Sales
	searchPattern := "%" + keyword + "%"

	err := r.db.WithContext(ctx).
		Where("code ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ? OR nick_name ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern).
		Order("code ASC").
		Offset(offset).
		Limit(limit).
		Find(&sales).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search sales by keyword: %w", err)
	}

	return sales, nil
}

func (r *SalesRepositoryImpl) GetWithBookingCounts(ctx context.Context, offset, limit int) ([]*models.Sales, error) {
	var sales []*models.Sales

	err := r.db.WithContext(ctx).
		Preload("Bookings").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&sales).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get sales with booking counts: %w", err)
	}

	return sales, nil
}

func (r *SalesRepositoryImpl) UpdateActiveStatus(ctx context.Context, id uuid.UUID, active bool) error {
	result := r.db.WithContext(ctx).
		Model(&models.Sales{}).
		Where("id = ?", id).
		Update("active", active)

	if result.Error != nil {
		return fmt.Errorf("failed to update sales active status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("sales with ID %s not found", id)
	}

	return nil
}

func (r *SalesRepositoryImpl) Activate(ctx context.Context, id uuid.UUID) error {
	return r.UpdateActiveStatus(ctx, id, true)
}

func (r *SalesRepositoryImpl) Deactivate(ctx context.Context, id uuid.UUID) error {
	return r.UpdateActiveStatus(ctx, id, false)
}

// WithTx returns a new repository instance with the given transaction
func (r *SalesRepositoryImpl) WithTx(tx *gorm.DB) repositories.SalesRepository {
	return &SalesRepositoryImpl{db: tx}
}