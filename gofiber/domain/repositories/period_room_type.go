package repositories

import (
	"context"
	"gofiber-template/domain/models"
	"github.com/google/uuid"
)

type PeriodRoomTypeRepository interface {
	Create(ctx context.Context, item *models.PeriodRoomType) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.PeriodRoomType, error)
	GetByPeriodID(ctx context.Context, periodID uuid.UUID) ([]*models.PeriodRoomType, error)
	Update(ctx context.Context, id uuid.UUID, item *models.PeriodRoomType) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByPeriodID(ctx context.Context, periodID uuid.UUID) error
	CreateBatch(ctx context.Context, items []*models.PeriodRoomType) error
	UpdateBatch(ctx context.Context, periodID uuid.UUID, items []*models.PeriodRoomType) error
}