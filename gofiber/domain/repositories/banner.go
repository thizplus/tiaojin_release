package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type BannerRepository interface {
	Create(ctx context.Context, entity *models.Banner) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Banner, error)
	Update(ctx context.Context, id uuid.UUID, entity *models.Banner) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.Banner, error)
	Count(ctx context.Context) (int64, error)

	SearchByTitle(ctx context.Context, title string, offset, limit int) ([]*models.Banner, error)
	GetPublished(ctx context.Context, offset, limit int) ([]*models.Banner, error)
	GetPublishedCount(ctx context.Context) (int64, error)
	GetRecentPublished(ctx context.Context, limit int) ([]*models.Banner, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.BannerStatus) error
	Publish(ctx context.Context, id uuid.UUID) error
	Unpublish(ctx context.Context, id uuid.UUID) error
}