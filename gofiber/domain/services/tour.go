package services

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
)


type TourService interface {
	CreateTour(ctx context.Context, req dto.CreateTourRequest) (*dto.TourResponse, error)
	GetTourByID(ctx context.Context, tourID uuid.UUID) (*dto.TourResponse, error)
	GetTourByCode(ctx context.Context, code string) (*dto.TourResponse, error)
	GetTourWithPeriods(ctx context.Context, tourID uuid.UUID) (*dto.TourWithPeriodsResponse, error)
	UpdateTour(ctx context.Context, tourID uuid.UUID, req dto.UpdateTourRequest) (*dto.TourResponse, error)
	DeleteTour(ctx context.Context, tourID uuid.UUID) error
	ListTours(ctx context.Context, req dto.TourListRequest) (*dto.TourListResponse, error)
	ListToursWithPeriods(ctx context.Context, req dto.TourListRequest, availableOnly bool) (*dto.TourWithPeriodsListResponse, error)
	ListSpecialTours(ctx context.Context, req dto.TourListRequest, availableOnly bool) (*dto.TourWithPeriodsListResponse, error)
	GetAvailableTours(ctx context.Context, offset, limit int) ([]*dto.TourWithPeriodsResponse, error)
	GetToursByCity(ctx context.Context, cityID uuid.UUID, offset, limit int) ([]*dto.TourResponse, error)
	GetPublishedTours(ctx context.Context, offset, limit int) ([]*dto.TourResponse, error)
	SearchTours(ctx context.Context, query string, offset, limit int) ([]*dto.TourResponse, error)
	UploadTourFile(ctx context.Context, req dto.UploadTourFileRequest, fileType string) (*dto.TourResponse, error)
	DeleteTourFile(ctx context.Context, tourID uuid.UUID, fileType string) (*dto.TourResponse, error)
	PublishTour(ctx context.Context, tourID uuid.UUID) (*dto.TourResponse, error)
	UnpublishTour(ctx context.Context, tourID uuid.UUID) (*dto.TourResponse, error)
	UpdateTourStatus(ctx context.Context, req dto.UpdateTourStatusRequest) (*dto.TourResponse, error)
	GetTourPeriods(ctx context.Context, tourID uuid.UUID) ([]*dto.PeriodResponse, error)
	GetTourPeriodsWithPagination(ctx context.Context, tourID uuid.UUID, offset, limit int) (*dto.PeriodListResponse, error)
	GetTourStats(ctx context.Context) (*dto.TourStatsResponse, error)
	ValidateTourCode(ctx context.Context, code string, excludeTourID *uuid.UUID) (bool, error)
}