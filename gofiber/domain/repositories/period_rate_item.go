package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type PeriodRateItemRepository interface {
	Create(ctx context.Context, item *models.PeriodRateItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.PeriodRateItem, error)
	GetByPeriodID(ctx context.Context, periodID uuid.UUID) ([]*models.PeriodRateItem, error)
	Update(ctx context.Context, id uuid.UUID, item *models.PeriodRateItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByPeriodID(ctx context.Context, periodID uuid.UUID) error
	CreateBatch(ctx context.Context, items []*models.PeriodRateItem) error
	UpdateBatch(ctx context.Context, periodID uuid.UUID, items []*models.PeriodRateItem) error
}