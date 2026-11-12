package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type BookingAdditionalServiceRepository interface {
	Create(ctx context.Context, entity *models.BookingAdditionalService) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.BookingAdditionalService, error)
	Update(ctx context.Context, id uuid.UUID, entity *models.BookingAdditionalService) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.BookingAdditionalService, error)
	Count(ctx context.Context) (int64, error)

	GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*models.BookingAdditionalService, error)
	DeleteByBookingID(ctx context.Context, bookingID uuid.UUID) error
	CreateBatch(ctx context.Context, entities []*models.BookingAdditionalService) error
	GetTotalByBookingID(ctx context.Context, bookingID uuid.UUID) (float64, error)
	GetTotalByBookingIDAndPassengerCount(ctx context.Context, bookingID uuid.UUID, passengerCount int) (float64, error)
}