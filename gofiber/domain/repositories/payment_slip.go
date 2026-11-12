package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type PaymentSlipRepository interface {
	Create(ctx context.Context, entity *models.PaymentSlip) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.PaymentSlip, error)
	Update(ctx context.Context, id uuid.UUID, entity *models.PaymentSlip) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.PaymentSlip, error)
	Count(ctx context.Context) (int64, error)

	GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*models.PaymentSlip, error)
	GetBySlipID(ctx context.Context, slipID string) (*models.PaymentSlip, error)
	GetByStatus(ctx context.Context, status models.PaymentSlipStatus, offset, limit int) ([]*models.PaymentSlip, error)
	GetPendingSlips(ctx context.Context, offset, limit int) ([]*models.PaymentSlip, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.PaymentSlipStatus) error
	GetTotalVerifiedAmountByBookingID(ctx context.Context, bookingID uuid.UUID) (float64, error)
	GetTotalActiveAmountByBookingID(ctx context.Context, bookingID uuid.UUID) (float64, error)
	GenerateSlipID(ctx context.Context) (string, error)

	// Reversal system methods
	ReversePaymentSlip(ctx context.Context, id uuid.UUID, reversedBy uuid.UUID, reason string) error
	GetPaymentSlipHistory(ctx context.Context, id uuid.UUID) (*models.PaymentSlip, error)
}