package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type SalesRepository interface {
	Create(ctx context.Context, entity *models.Sales) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Sales, error)
	Update(ctx context.Context, id uuid.UUID, entity *models.Sales) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.Sales, error)
	Count(ctx context.Context) (int64, error)

	GetByCode(ctx context.Context, code string) (*models.Sales, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Sales, error)
	GetActive(ctx context.Context, offset, limit int) ([]*models.Sales, error)
	SearchByName(ctx context.Context, name string, offset, limit int) ([]*models.Sales, error)
	SearchByKeyword(ctx context.Context, keyword string, offset, limit int) ([]*models.Sales, error)
	GetWithBookingCounts(ctx context.Context, offset, limit int) ([]*models.Sales, error)
	UpdateActiveStatus(ctx context.Context, id uuid.UUID, active bool) error
	Activate(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id uuid.UUID) error
}