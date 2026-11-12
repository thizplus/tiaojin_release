package repositories

import (
	"context"
	"time"

	"gofiber-template/domain/models"

	"github.com/google/uuid"
)

// SeatHoldRepository defines the interface for seat hold operations
type SeatHoldRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, seatHold *models.SeatHold) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.SeatHold, error)
	Update(ctx context.Context, seatHold *models.SeatHold) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetByPeriodID(ctx context.Context, periodID uuid.UUID, includeInactive bool) ([]*models.SeatHold, error)
	GetByAdminID(ctx context.Context, adminID uuid.UUID, offset, limit int) ([]*models.SeatHold, error)
	GetActiveHolds(ctx context.Context, offset, limit int) ([]*models.SeatHold, error)
	GetExpiredHolds(ctx context.Context, beforeTime time.Time) ([]*models.SeatHold, error)

	// Aggregation operations
	GetTotalHeldSeats(ctx context.Context, periodID uuid.UUID) (int, error)
	GetActiveHeldSeats(ctx context.Context, periodID uuid.UUID) (int, error)
	GetSeatHoldSummary(ctx context.Context, periodID uuid.UUID) (*models.SeatHoldSummary, error)

	// Batch operations
	ExpireOldHolds(ctx context.Context, beforeTime time.Time) (int, error)
	ReleaseHoldsByAdmin(ctx context.Context, adminID uuid.UUID) (int, error)
	ReleaseHoldsByPeriod(ctx context.Context, periodID uuid.UUID) (int, error)

	// Validation operations
	CheckAvailableSeatsForHold(ctx context.Context, periodID uuid.UUID, requestedSeats int) (bool, error)
	GetConflictingHolds(ctx context.Context, periodID uuid.UUID, adminID uuid.UUID) ([]*models.SeatHold, error)

	// Advanced queries
	GetHoldsExpiringBefore(ctx context.Context, beforeTime time.Time, limit int) ([]*models.SeatHold, error)
	GetHoldsByStatus(ctx context.Context, status models.SeatHoldStatus, offset, limit int) ([]*models.SeatHold, error)
	GetAdminHoldStatistics(ctx context.Context, adminID uuid.UUID) (map[string]interface{}, error)

	// Count operations
	CountByPeriodID(ctx context.Context, periodID uuid.UUID) (int64, error)
	CountByAdminID(ctx context.Context, adminID uuid.UUID) (int64, error)
	CountByStatus(ctx context.Context, status models.SeatHoldStatus) (int64, error)

	// Seat availability calculations
	FindActiveBYPeriodID(periodID uuid.UUID) ([]*models.SeatHold, error)
}