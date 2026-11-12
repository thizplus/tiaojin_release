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

type TourRepositoryImpl struct {
	db *gorm.DB
}

func NewTourRepository(db *gorm.DB) repositories.TourRepository {
	return &TourRepositoryImpl{db: db}
}

func (r *TourRepositoryImpl) Create(ctx context.Context, entity *models.Tour) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create tour: %w", err)
	}
	return nil
}

func (r *TourRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Tour, error) {
	var tour models.Tour
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tour).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get tour by ID: %w", err)
	}

	return &tour, nil
}

func (r *TourRepositoryImpl) Update(ctx context.Context, id uuid.UUID, entity *models.Tour) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(entity)

	if result.Error != nil {
		return fmt.Errorf("failed to update tour: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("tour with ID %s not found", id)
	}

	return nil
}

func (r *TourRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Tour{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete tour: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("tour with ID %s not found", id)
	}

	return nil
}

func (r *TourRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Tour, error) {
	var tours []*models.Tour

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&tours).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list tours: %w", err)
	}

	return tours, nil
}

func (r *TourRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Tour{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count tours: %w", err)
	}

	return count, nil
}

func (r *TourRepositoryImpl) GetByCode(ctx context.Context, code string) (*models.Tour, error) {
	var tour models.Tour
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&tour).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get tour by code: %w", err)
	}

	return &tour, nil
}

func (r *TourRepositoryImpl) GetByCity(ctx context.Context, cityID uuid.UUID, offset, limit int) ([]*models.Tour, error) {
	var tours []*models.Tour

	err := r.db.WithContext(ctx).
		Where("city_id = ?", cityID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&tours).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get tours by city: %w", err)
	}

	return tours, nil
}

func (r *TourRepositoryImpl) GetPublished(ctx context.Context, offset, limit int) ([]*models.Tour, error) {
	var tours []*models.Tour

	err := r.db.WithContext(ctx).
		Where("status = ?", models.TourStatusPublish).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&tours).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get published tours: %w", err)
	}

	return tours, nil
}

func (r *TourRepositoryImpl) GetPublishedByCity(ctx context.Context, cityID uuid.UUID, offset, limit int) ([]*models.Tour, error) {
	var tours []*models.Tour

	err := r.db.WithContext(ctx).
		Where("city_id = ? AND status = ?", cityID, models.TourStatusPublish).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&tours).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get published tours by city: %w", err)
	}

	return tours, nil
}

func (r *TourRepositoryImpl) SearchByTitle(ctx context.Context, title string, offset, limit int) ([]*models.Tour, error) {
	var tours []*models.Tour
	searchPattern := "%" + title + "%"

	err := r.db.WithContext(ctx).
		Where("title ILIKE ?", searchPattern).
		Order("title ASC").
		Offset(offset).
		Limit(limit).
		Find(&tours).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search tours by title: %w", err)
	}

	return tours, nil
}

func (r *TourRepositoryImpl) SearchByTitleOrCode(ctx context.Context, search string, offset, limit int) ([]*models.Tour, error) {
	var tours []*models.Tour
	searchPattern := "%" + search + "%"

	err := r.db.WithContext(ctx).
		Where("title ILIKE ? OR code ILIKE ?", searchPattern, searchPattern).
		Order("title ASC").
		Offset(offset).
		Limit(limit).
		Find(&tours).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search tours by title or code: %w", err)
	}

	return tours, nil
}

func (r *TourRepositoryImpl) GetWithPeriods(ctx context.Context, id uuid.UUID) (*models.Tour, error) {
	var tour models.Tour
	err := r.db.WithContext(ctx).
		Preload("Periods").
		Where("id = ?", id).
		First(&tour).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get tour with periods: %w", err)
	}

	return &tour, nil
}

func (r *TourRepositoryImpl) GetPublishedWithAvailablePeriods(ctx context.Context, offset, limit int) ([]*models.Tour, error) {
	var tours []*models.Tour

	err := r.db.WithContext(ctx).
		Preload("Periods", "status = ? AND start_date > NOW()", models.PeriodStatusPublish).
		Where("status = ?", models.TourStatusPublish).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&tours).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get published tours with available periods: %w", err)
	}

	return tours, nil
}

// GetToursWithSalePeriods ดึงเฉพาะทัวร์ที่มี period แบบ sale จาก DB ตรงๆ
func (r *TourRepositoryImpl) GetToursWithSalePeriods(ctx context.Context, offset, limit int) ([]*models.Tour, error) {
	var tours []*models.Tour

	subQuery := r.db.WithContext(ctx).
		Model(&models.Period{}).
		Select("DISTINCT tour_id").
		Where("status = ? AND start_date > NOW() AND (rate_type = 'sale' OR commission_type = 'sale') AND deleted_at IS NULL", models.PeriodStatusPublish)

	err := r.db.WithContext(ctx).
		Preload("Periods", "status = ? AND start_date > NOW() AND (rate_type = 'sale' OR commission_type = 'sale') AND deleted_at IS NULL", models.PeriodStatusPublish).
		Where("status = ? AND id IN (?)", models.TourStatusPublish, subQuery).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&tours).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get tours with sale periods: %w", err)
	}

	return tours, nil
}

func (r *TourRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TourStatus) error {
	result := r.db.WithContext(ctx).
		Model(&models.Tour{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update tour status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("tour with ID %s not found", id)
	}

	return nil
}

// WithTx returns a new repository instance with the given transaction
func (r *TourRepositoryImpl) WithTx(tx *gorm.DB) repositories.TourRepository {
	return &TourRepositoryImpl{db: tx}
}