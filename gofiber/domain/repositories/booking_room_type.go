package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type BookingRoomTypeRepository interface {
	Create(ctx context.Context, entity *models.BookingRoomType) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.BookingRoomType, error)
	Update(ctx context.Context, id uuid.UUID, entity *models.BookingRoomType) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.BookingRoomType, error)
	Count(ctx context.Context) (int64, error)

	GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*models.BookingRoomType, error)
	DeleteByBookingID(ctx context.Context, bookingID uuid.UUID) error
	CreateBatch(ctx context.Context, entities []*models.BookingRoomType) error
	GetTotalByBookingID(ctx context.Context, bookingID uuid.UUID) (float64, error)
}