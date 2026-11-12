package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
)

type PeriodRepository interface {
	Create(ctx context.Context, entity *models.Period) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Period, error)
	Update(ctx context.Context, id uuid.UUID, entity *models.Period) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.Period, error)
	Count(ctx context.Context) (int64, error)

	GetByTour(ctx context.Context, tourID uuid.UUID, offset, limit int) ([]*models.Period, error)
	GetByTourID(ctx context.Context, tourID uuid.UUID) ([]*models.Period, error)
	GetByTourIDWithPagination(ctx context.Context, tourID uuid.UUID, offset, limit int) ([]*models.Period, error)
	CountByTourID(ctx context.Context, tourID uuid.UUID) (int64, error)
	GetAvailable(ctx context.Context, offset, limit int) ([]*models.Period, error)
	GetAvailableByTour(ctx context.Context, tourID uuid.UUID, offset, limit int) ([]*models.Period, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Period, error)
	GetUpcoming(ctx context.Context, offset, limit int) ([]*models.Period, error)
	UpdateBookingCount(ctx context.Context, id uuid.UUID, bookingCount int) error
	IncrementBookingCount(ctx context.Context, id uuid.UUID, increment int) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.PeriodStatus) error
	ListWithFilters(ctx context.Context, req dto.PeriodListRequest) ([]*models.Period, int64, error)
}