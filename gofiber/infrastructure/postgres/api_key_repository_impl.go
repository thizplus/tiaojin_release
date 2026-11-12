package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type APIKeyRepositoryImpl struct {
	db *gorm.DB
}

func NewAPIKeyRepository(db *gorm.DB) repositories.APIKeyRepository {
	return &APIKeyRepositoryImpl{db: db}
}

func (r *APIKeyRepositoryImpl) Create(ctx context.Context, apiKey *models.APIKey) error {
	if err := r.db.WithContext(ctx).Create(apiKey).Error; err != nil {
		return fmt.Errorf("failed to create API key: %w", err)
	}
	return nil
}

func (r *APIKeyRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.APIKey, error) {
	var apiKey models.APIKey
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", id).
		First(&apiKey).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get API key by ID: %w", err)
	}

	return &apiKey, nil
}

func (r *APIKeyRepositoryImpl) GetByKeyHash(ctx context.Context, keyHash string) (*models.APIKey, error) {
	var apiKey models.APIKey
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("key_hash = ?", keyHash).
		First(&apiKey).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get API key by hash: %w", err)
	}

	return &apiKey, nil
}

func (r *APIKeyRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.APIKey, error) {
	var apiKeys []*models.APIKey
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&apiKeys).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get API keys by user ID: %w", err)
	}

	return apiKeys, nil
}

func (r *APIKeyRepositoryImpl) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.APIKey{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count API keys: %w", err)
	}

	return count, nil
}

func (r *APIKeyRepositoryImpl) ListWithFilters(ctx context.Context, req dto.ListAPIKeysRequest, userID uuid.UUID) ([]*models.APIKey, error) {
	var apiKeys []*models.APIKey
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	// Filter by active status
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	// Filter by search term (name or description)
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Filter expired keys
	if !req.ShowExpired {
		query = query.Where("expires_at IS NULL OR expires_at > ?", time.Now())
	}

	err := query.
		Order("created_at DESC").
		Offset(req.Offset).
		Limit(req.Limit).
		Find(&apiKeys).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list API keys with filters: %w", err)
	}

	return apiKeys, nil
}

func (r *APIKeyRepositoryImpl) CountWithFilters(ctx context.Context, req dto.ListAPIKeysRequest, userID uuid.UUID) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.APIKey{}).Where("user_id = ?", userID)

	// Filter by active status
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	// Filter by search term
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Filter expired keys
	if !req.ShowExpired {
		query = query.Where("expires_at IS NULL OR expires_at > ?", time.Now())
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count API keys with filters: %w", err)
	}

	return count, nil
}

func (r *APIKeyRepositoryImpl) Update(ctx context.Context, id uuid.UUID, apiKey *models.APIKey) error {
	// Use Select to explicitly update all fields including zero values like is_active=false
	result := r.db.WithContext(ctx).
		Model(&models.APIKey{}).
		Where("id = ?", id).
		Select("name", "description", "is_active", "expires_at", "ip_whitelist", "permissions", "updated_at").
		Updates(apiKey)

	if result.Error != nil {
		return fmt.Errorf("failed to update API key: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

func (r *APIKeyRepositoryImpl) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&models.APIKey{}).
		Where("id = ?", id).
		Update("last_used_at", now)

	if result.Error != nil {
		return fmt.Errorf("failed to update last used timestamp: %w", result.Error)
	}

	return nil
}

func (r *APIKeyRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.APIKey{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete API key: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

func (r *APIKeyRepositoryImpl) SoftDelete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.APIKey{})

	if result.Error != nil {
		return fmt.Errorf("failed to soft delete API key: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}