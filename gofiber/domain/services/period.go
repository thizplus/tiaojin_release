package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
)


type PeriodService interface {
	CreatePeriod(ctx context.Context, req dto.CreatePeriodRequest) (*dto.PeriodResponse, error)
	GetPeriodByID(ctx context.Context, periodID uuid.UUID) (*dto.PeriodResponse, error)
	GetPeriodWithTour(ctx context.Context, periodID uuid.UUID) (*dto.PeriodResponse, error)
	UpdatePeriod(ctx context.Context, periodID uuid.UUID, req dto.UpdatePeriodRequest) (*dto.PeriodResponse, error)
	DeletePeriod(ctx context.Context, periodID uuid.UUID) error
	ListPeriods(ctx context.Context, req dto.PeriodListRequest) (*dto.PeriodListResponse, error)
	GetPeriodsByTour(ctx context.Context, tourID uuid.UUID, offset, limit int) ([]*dto.PeriodResponse, error)
	GetAvailablePeriods(ctx context.Context, offset, limit int) ([]*dto.PeriodResponse, error)
	GetUpcomingPeriods(ctx context.Context, days int, offset, limit int) ([]*dto.PeriodResponse, error)

	// Availability Management
	CheckAvailability(ctx context.Context, periodID uuid.UUID, requestedSeats int) (*dto.AvailabilityCheckResponse, error)
	ReserveSeats(ctx context.Context, periodID uuid.UUID, seats int) error
	ReleaseSeats(ctx context.Context, periodID uuid.UUID, seats int) error

	// File Management
	UploadPeriodFile(ctx context.Context, periodID uuid.UUID, fileData []byte, fileName, fileType string) (*dto.PeriodResponse, error)
	DeletePeriodFile(ctx context.Context, periodID uuid.UUID, fileURL string) (*dto.PeriodResponse, error)

	// Business Logic
	PublishPeriod(ctx context.Context, periodID uuid.UUID) (*dto.PeriodResponse, error)
	UnpublishPeriod(ctx context.Context, periodID uuid.UUID) (*dto.PeriodResponse, error)
	UpdatePeriodStatus(ctx context.Context, req dto.UpdatePeriodStatusRequest) (*dto.PeriodResponse, error)
	GetPeriodAvailability(ctx context.Context, periodID uuid.UUID) (*dto.PeriodAvailabilityResponse, error)
	CalculateRealTimeBookingStats(ctx context.Context, periodID uuid.UUID) (*dto.PeriodBookingStatsResponse, error)
	GetPeriodStats(ctx context.Context) (*dto.PeriodStatsResponse, error)
	ValidatePeriodDates(ctx context.Context, startDate, endDate time.Time) error
	CalculateCommission(ctx context.Context, periodID uuid.UUID, rateType string) (float64, error)
	CalculateCommissionForBooking(ctx context.Context, periodID uuid.UUID, rateType string, passengers []models.RateItem) (float64, error)
}