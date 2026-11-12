package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type CityRepository interface {
	Create(ctx context.Context, entity *models.City) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.City, error)
	Update(ctx context.Context, id uuid.UUID, entity *models.City) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.City, error)
	Count(ctx context.Context) (int64, error)

	GetPublished(ctx context.Context, offset, limit int) ([]*models.City, error)
	SearchByTitle(ctx context.Context, title string, offset, limit int) ([]*models.City, error)
	GetWithTours(ctx context.Context, id uuid.UUID) (*models.City, error)
	GetPublishedWithActiveTours(ctx context.Context, offset, limit int) ([]*models.City, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.CityStatus) error
}