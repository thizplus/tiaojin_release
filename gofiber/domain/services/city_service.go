package services

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
)

type CityService interface {
	// CRUD Operations
	CreateCity(ctx context.Context, req dto.CreateCityRequest) (*dto.CityResponse, error)
	UpdateCity(ctx context.Context, cityID uuid.UUID, req dto.UpdateCityRequest) (*dto.CityResponse, error)
	DeleteCity(ctx context.Context, cityID uuid.UUID) error
	GetCityByID(ctx context.Context, cityID uuid.UUID) (*dto.CityResponse, error)
	ListCities(ctx context.Context, req dto.CityListRequest) (*dto.CityListResponse, error)

	// Status Management
	UpdateStatus(ctx context.Context, req dto.UpdateCityStatusRequest) (*dto.CityResponse, error)
	UpdateCityStatus(ctx context.Context, cityID uuid.UUID, req dto.UpdateCityStatusRequest) (*dto.CityResponse, error)

	// Public Operations
	GetPublishedCities(ctx context.Context, offset, limit int) ([]*dto.CityResponse, error)
	GetCityTours(ctx context.Context, cityID uuid.UUID, offset, limit int) ([]*dto.TourResponse, error)
}