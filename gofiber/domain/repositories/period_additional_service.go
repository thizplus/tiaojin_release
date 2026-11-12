package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type PeriodAdditionalServiceRepository interface {
	Create(ctx context.Context, item *models.PeriodAdditionalService) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.PeriodAdditionalService, error)
	GetByPeriodID(ctx context.Context, periodID uuid.UUID) ([]*models.PeriodAdditionalService, error)
	Update(ctx context.Context, id uuid.UUID, item *models.PeriodAdditionalService) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByPeriodID(ctx context.Context, periodID uuid.UUID) error
	CreateBatch(ctx context.Context, items []*models.PeriodAdditionalService) error
	UpdateBatch(ctx context.Context, periodID uuid.UUID, items []*models.PeriodAdditionalService) error
}