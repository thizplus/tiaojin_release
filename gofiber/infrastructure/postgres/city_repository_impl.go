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

type CityRepositoryImpl struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) repositories.CityRepository {
	return &CityRepositoryImpl{db: db}
}

func (r *CityRepositoryImpl) Create(ctx context.Context, entity *models.City) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create city: %w", err)
	}
	return nil
}

func (r *CityRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.City, error) {
	var city models.City
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&city).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get city by ID: %w", err)
	}

	return &city, nil
}

func (r *CityRepositoryImpl) Update(ctx context.Context, id uuid.UUID, entity *models.City) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(entity)

	if result.Error != nil {
		return fmt.Errorf("failed to update city: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("city with ID %s not found", id)
	}

	return nil
}

func (r *CityRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.City{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete city: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("city with ID %s not found", id)
	}

	return nil
}

func (r *CityRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.City, error) {
	var cities []*models.City

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&cities).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list cities: %w", err)
	}

	return cities, nil
}

func (r *CityRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.City{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count cities: %w", err)
	}

	return count, nil
}

func (r *CityRepositoryImpl) GetPublished(ctx context.Context, offset, limit int) ([]*models.City, error) {
	var cities []*models.City

	err := r.db.WithContext(ctx).
		Where("status = ?", models.CityStatusPublish).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&cities).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get published cities: %w", err)
	}

	return cities, nil
}

func (r *CityRepositoryImpl) SearchByTitle(ctx context.Context, title string, offset, limit int) ([]*models.City, error) {
	var cities []*models.City
	searchPattern := "%" + title + "%"

	err := r.db.WithContext(ctx).
		Where("title ILIKE ?", searchPattern).
		Order("title ASC").
		Offset(offset).
		Limit(limit).
		Find(&cities).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search cities by title: %w", err)
	}

	return cities, nil
}

func (r *CityRepositoryImpl) GetWithTours(ctx context.Context, id uuid.UUID) (*models.City, error) {
	var city models.City
	err := r.db.WithContext(ctx).
		Preload("Tours").
		Where("id = ?", id).
		First(&city).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get city with tours: %w", err)
	}

	return &city, nil
}

func (r *CityRepositoryImpl) GetPublishedWithActiveTours(ctx context.Context, offset, limit int) ([]*models.City, error) {
	var cities []*models.City

	err := r.db.WithContext(ctx).
		Preload("Tours", "status = ?", models.TourStatusPublish).
		Where("status = ?", models.CityStatusPublish).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&cities).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get published cities with active tours: %w", err)
	}

	return cities, nil
}

func (r *CityRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status models.CityStatus) error {
	result := r.db.WithContext(ctx).
		Model(&models.City{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update city status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("city with ID %s not found", id)
	}

	return nil
}

// WithTx returns a new repository instance with the given transaction
func (r *CityRepositoryImpl) WithTx(tx *gorm.DB) repositories.CityRepository {
	return &CityRepositoryImpl{db: tx}
}