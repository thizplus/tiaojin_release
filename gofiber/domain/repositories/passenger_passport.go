package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type PassengerPassportRepository interface {
	Create(ctx context.Context, entity *models.PassengerPassport) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.PassengerPassport, error)
	Update(ctx context.Context, id uuid.UUID, entity *models.PassengerPassport) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.PassengerPassport, error)
	Count(ctx context.Context) (int64, error)

	GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*models.PassengerPassport, error)
	GetByPassportNumber(ctx context.Context, passportNumber string) (*models.PassengerPassport, error)
	GetByStatus(ctx context.Context, status models.PassportStatus, offset, limit int) ([]*models.PassengerPassport, error)
	GetPendingProcessing(ctx context.Context, offset, limit int) ([]*models.PassengerPassport, error)
	GetExpiring(ctx context.Context, daysBeforeExpiry int, offset, limit int) ([]*models.PassengerPassport, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.PassportStatus) error
	UpdatePersonalInfo(ctx context.Context, id uuid.UUID, passport *models.PassengerPassport) error
	CreateBatch(ctx context.Context, entities []*models.PassengerPassport) error
	DeleteByBookingID(ctx context.Context, bookingID uuid.UUID) error
}