package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
)

type BookingRepository interface {
	Create(ctx context.Context, entity *models.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Booking, error)
	GetByIDWithPaymentSlips(ctx context.Context, id uuid.UUID) (*models.Booking, error)
	Update(ctx context.Context, id uuid.UUID, entity *models.Booking) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.Booking, error)
	Count(ctx context.Context) (int64, error)
	ListWithFilters(ctx context.Context, req dto.BookingListRequest) ([]*models.Booking, error)
	CountWithFilters(ctx context.Context, req dto.BookingListRequest) (int64, error)

	GetByInvoiceID(ctx context.Context, invoiceID string) (*models.Booking, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Booking, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	GetByPeriodID(ctx context.Context, periodID uuid.UUID, offset, limit int) ([]*models.Booking, error)
	GetBySalesID(ctx context.Context, salesID uuid.UUID, offset, limit int) ([]*models.Booking, error)
	GetByStatus(ctx context.Context, status models.BookingStatus, offset, limit int) ([]*models.Booking, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Booking, error)
	GetWithDetails(ctx context.Context, id uuid.UUID) (*models.Booking, error)
	GetPendingPayments(ctx context.Context, offset, limit int) ([]*models.Booking, error)
	GetOverduePayments(ctx context.Context, offset, limit int) ([]*models.Booking, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.BookingStatus) error
	UpdatePaymentAmount(ctx context.Context, id uuid.UUID, amount float64) error
	GenerateInvoiceID(ctx context.Context, year int) (string, error)

	// Period booking stats
	CountBookingsByPeriodAndStatus(ctx context.Context, periodID uuid.UUID, status string) (int, int, error)

	// Seat availability calculations
	FindByPeriodIDAndStatuses(periodID uuid.UUID, statuses []string) ([]*models.Booking, error)
}