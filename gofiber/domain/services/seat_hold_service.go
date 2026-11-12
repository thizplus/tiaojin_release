package services

import (
	"context"
	"time"

	"gofiber-template/domain/dto"

	"github.com/google/uuid"
)

// SeatHoldService defines the interface for seat hold business operations
type SeatHoldService interface {
	// Core seat hold operations
	CreateSeatHold(ctx context.Context, req dto.CreateSeatHoldRequest) (*dto.SeatHoldResponse, error)
	GetSeatHoldByID(ctx context.Context, id uuid.UUID) (*dto.SeatHoldResponse, error)
	UpdateSeatHold(ctx context.Context, id uuid.UUID, req dto.UpdateSeatHoldRequest) (*dto.SeatHoldResponse, error)
	ReleaseSeatHold(ctx context.Context, id uuid.UUID, reason string) (*dto.SeatHoldResponse, error)
	DeleteSeatHold(ctx context.Context, id uuid.UUID) error

	// Query operations
	GetSeatHoldsByPeriod(ctx context.Context, periodID uuid.UUID, includeInactive bool) ([]*dto.SeatHoldResponse, error)
	GetSeatHoldsByAdmin(ctx context.Context, adminID uuid.UUID, offset, limit int) (*dto.SeatHoldListResponse, error)
	GetActiveSeatHolds(ctx context.Context, offset, limit int) (*dto.SeatHoldListResponse, error)

	// Business operations
	ExtendSeatHold(ctx context.Context, id uuid.UUID, duration time.Duration) (*dto.SeatHoldResponse, error)
	ConvertSeatHoldToBooking(ctx context.Context, id uuid.UUID, bookingReq dto.CreateBookingRequest) (*dto.SeatHoldResponse, error)
	CheckSeatAvailability(ctx context.Context, periodID uuid.UUID, requestedSeats int) (*dto.SeatAvailabilityResponse, error)

	// Batch operations
	ExpireOldSeatHolds(ctx context.Context) (*dto.ExpireHoldsResponse, error)
	ReleaseAdminSeatHolds(ctx context.Context, adminID uuid.UUID, reason string) (*dto.ReleaseHoldsResponse, error)
	ReleasePeriodSeatHolds(ctx context.Context, periodID uuid.UUID, reason string) (*dto.ReleaseHoldsResponse, error)

	// Analytics and reporting
	GetSeatHoldSummary(ctx context.Context, periodID uuid.UUID) (*dto.SeatHoldSummaryResponse, error)
	GetAdminSeatHoldStatistics(ctx context.Context, adminID uuid.UUID) (*dto.AdminSeatHoldStatsResponse, error)
	GetSeatHoldAnalytics(ctx context.Context, req dto.SeatHoldAnalyticsRequest) (*dto.SeatHoldAnalyticsResponse, error)

	// Validation and utility
	ValidateSeatHoldRequest(ctx context.Context, req dto.CreateSeatHoldRequest) error
	CalculateOptimalExpiry(ctx context.Context, periodID uuid.UUID) (time.Time, error)
	GetUpcomingExpirations(ctx context.Context, hoursAhead int, limit int) ([]*dto.SeatHoldResponse, error)

	// Dynamic Job-based Expiration Management
	CreateSeatHoldExpiryJob(ctx context.Context, seatHoldID uuid.UUID, expiresAt time.Time) error
	DeleteSeatHoldExpiryJob(ctx context.Context, seatHoldID uuid.UUID) error
	ExpireSeatHold(ctx context.Context, seatHoldID uuid.UUID) error
}