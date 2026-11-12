package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
)

type APIKeyRepository interface {
	// Basic CRUD
	Create(ctx context.Context, apiKey *models.APIKey) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.APIKey, error)
	Update(ctx context.Context, id uuid.UUID, apiKey *models.APIKey) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query methods
	GetByKeyHash(ctx context.Context, keyHash string) (*models.APIKey, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.APIKey, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)

	// List with filters
	ListWithFilters(ctx context.Context, req dto.ListAPIKeysRequest, userID uuid.UUID) ([]*models.APIKey, error)
	CountWithFilters(ctx context.Context, req dto.ListAPIKeysRequest, userID uuid.UUID) (int64, error)

	// Update last used timestamp
	UpdateLastUsed(ctx context.Context, id uuid.UUID) error

	// Soft delete
	SoftDelete(ctx context.Context, id uuid.UUID) error
}