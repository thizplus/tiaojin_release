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

type BannerRepositoryImpl struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) repositories.BannerRepository {
	return &BannerRepositoryImpl{db: db}
}

func (r *BannerRepositoryImpl) Create(ctx context.Context, entity *models.Banner) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create banner: %w", err)
	}
	return nil
}

func (r *BannerRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Banner, error) {
	var banner models.Banner
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&banner).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get banner by ID: %w", err)
	}

	return &banner, nil
}

func (r *BannerRepositoryImpl) Update(ctx context.Context, id uuid.UUID, entity *models.Banner) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(entity)

	if result.Error != nil {
		return fmt.Errorf("failed to update banner: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("banner with ID %s not found", id)
	}

	return nil
}

func (r *BannerRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Banner{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete banner: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("banner with ID %s not found", id)
	}

	return nil
}

func (r *BannerRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Banner, error) {
	var banners []*models.Banner

	err := r.db.WithContext(ctx).
		Order("display_order ASC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&banners).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list banners: %w", err)
	}

	return banners, nil
}

func (r *BannerRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Banner{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count banners: %w", err)
	}

	return count, nil
}

func (r *BannerRepositoryImpl) SearchByTitle(ctx context.Context, title string, offset, limit int) ([]*models.Banner, error) {
	var banners []*models.Banner
	searchPattern := "%" + title + "%"

	err := r.db.WithContext(ctx).
		Where("title ILIKE ?", searchPattern).
		Order("title ASC").
		Offset(offset).
		Limit(limit).
		Find(&banners).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search banners by title: %w", err)
	}

	return banners, nil
}

func (r *BannerRepositoryImpl) GetPublished(ctx context.Context, offset, limit int) ([]*models.Banner, error) {
	var banners []*models.Banner

	err := r.db.WithContext(ctx).
		Where("status = ?", models.BannerStatusPublish).
		Order("display_order ASC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&banners).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get published banners: %w", err)
	}

	return banners, nil
}

func (r *BannerRepositoryImpl) GetPublishedCount(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.Banner{}).
		Where("status = ?", models.BannerStatusPublish).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count published banners: %w", err)
	}

	return count, nil
}

func (r *BannerRepositoryImpl) GetRecentPublished(ctx context.Context, limit int) ([]*models.Banner, error) {
	var banners []*models.Banner

	err := r.db.WithContext(ctx).
		Where("status = ?", models.BannerStatusPublish).
		Order("display_order ASC, created_at DESC").
		Limit(limit).
		Find(&banners).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get recent published banners: %w", err)
	}

	return banners, nil
}

func (r *BannerRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status models.BannerStatus) error {
	result := r.db.WithContext(ctx).
		Model(&models.Banner{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update banner status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("banner with ID %s not found", id)
	}

	return nil
}

func (r *BannerRepositoryImpl) Publish(ctx context.Context, id uuid.UUID) error {
	return r.UpdateStatus(ctx, id, models.BannerStatusPublish)
}

func (r *BannerRepositoryImpl) Unpublish(ctx context.Context, id uuid.UUID) error {
	return r.UpdateStatus(ctx, id, models.BannerStatusDraft)
}

// WithTx returns a new repository instance with the given transaction
func (r *BannerRepositoryImpl) WithTx(tx *gorm.DB) repositories.BannerRepository {
	return &BannerRepositoryImpl{db: tx}
}