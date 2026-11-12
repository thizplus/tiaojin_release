package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type TourRepository interface {
	Create(ctx context.Context, entity *models.Tour) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Tour, error)
	Update(ctx context.Context, id uuid.UUID, entity *models.Tour) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.Tour, error)
	Count(ctx context.Context) (int64, error)

	GetByCode(ctx context.Context, code string) (*models.Tour, error)
	GetByCity(ctx context.Context, cityID uuid.UUID, offset, limit int) ([]*models.Tour, error)
	GetPublished(ctx context.Context, offset, limit int) ([]*models.Tour, error)
	GetPublishedByCity(ctx context.Context, cityID uuid.UUID, offset, limit int) ([]*models.Tour, error)
	SearchByTitle(ctx context.Context, title string, offset, limit int) ([]*models.Tour, error)
	SearchByTitleOrCode(ctx context.Context, search string, offset, limit int) ([]*models.Tour, error)
	GetWithPeriods(ctx context.Context, id uuid.UUID) (*models.Tour, error)
	GetPublishedWithAvailablePeriods(ctx context.Context, offset, limit int) ([]*models.Tour, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.TourStatus) error
	GetToursWithSalePeriods(ctx context.Context, offset, limit int) ([]*models.Tour, error)
}