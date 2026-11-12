package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type BookingRateTypeRepository interface {
	Create(ctx context.Context, entity *models.BookingRateType) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.BookingRateType, error)
	Update(ctx context.Context, id uuid.UUID, entity *models.BookingRateType) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.BookingRateType, error)
	Count(ctx context.Context) (int64, error)

	GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*models.BookingRateType, error)
	DeleteByBookingID(ctx context.Context, bookingID uuid.UUID) error
	CreateBatch(ctx context.Context, entities []*models.BookingRateType) error
	GetTotalByBookingID(ctx context.Context, bookingID uuid.UUID) (float64, error)
	GetIncludedByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*models.BookingRateType, error)
	GetExcludedByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*models.BookingRateType, error)
}