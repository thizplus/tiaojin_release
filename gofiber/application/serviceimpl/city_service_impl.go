package serviceimpl

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/cdn"
	apperrors "gofiber-template/pkg/errors"
	"gofiber-template/pkg/utils"
)

type CityServiceImpl struct {
	cityRepo   repositories.CityRepository
	tourRepo   repositories.TourRepository
	periodRepo repositories.PeriodRepository
}

func NewCityService(cityRepo repositories.CityRepository, tourRepo repositories.TourRepository, periodRepo repositories.PeriodRepository) services.CityService {
	return &CityServiceImpl{
		cityRepo:   cityRepo,
		tourRepo:   tourRepo,
		periodRepo: periodRepo,
	}
}

// CRUD Operations
func (s *CityServiceImpl) CreateCity(ctx context.Context, req dto.CreateCityRequest) (*dto.CityResponse, error) {
	// Create domain model
	city := &models.City{
		Title:     req.Title,
		Detail:    req.Detail,
		Image:     req.Image,
		Status:    models.CityStatus(req.Status), // Use status from request
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create city
	if err := s.cityRepo.Create(ctx, city); err != nil {
		log.Printf("Error creating city: %v", err)
		return nil, fmt.Errorf("failed to create city: %w", err)
	}

	return s.toCityResponse(ctx, city, nil), nil
}

func (s *CityServiceImpl) UpdateCity(ctx context.Context, cityID uuid.UUID, req dto.UpdateCityRequest) (*dto.CityResponse, error) {
	// Get existing city
	city, err := s.cityRepo.GetByID(ctx, cityID)
	if err != nil {
		log.Printf("Error retrieving city: %v", err)
		return nil, fmt.Errorf("failed to get city: %w", err)
	}

	if city == nil {
		return nil, apperrors.ErrCityNotFound
	}

	// Update fields if provided
	if req.Title != nil {
		city.Title = *req.Title
	}
	if req.Detail != nil {
		city.Detail = req.Detail
	}
	if req.Image != nil {
		city.Image = *req.Image
	}
	if req.Status != nil {
		city.Status = models.CityStatus(*req.Status)
	}

	city.UpdatedAt = time.Now()

	// Update city
	if err := s.cityRepo.Update(ctx, cityID, city); err != nil {
		log.Printf("Error updating city: %v", err)
		return nil, fmt.Errorf("failed to update city: %w", err)
	}

	return s.toCityResponse(ctx, city, nil), nil
}

func (s *CityServiceImpl) DeleteCity(ctx context.Context, cityID uuid.UUID) error {
	// Check if city exists
	city, err := s.cityRepo.GetByID(ctx, cityID)
	if err != nil {
		log.Printf("Error retrieving city: %v", err)
		return fmt.Errorf("failed to get city: %w", err)
	}

	if city == nil {
		return apperrors.ErrCityNotFound
	}

	// Delete city
	if err := s.cityRepo.Delete(ctx, cityID); err != nil {
		log.Printf("Error deleting city: %v", err)
		return fmt.Errorf("failed to delete city: %w", err)
	}

	return nil
}

func (s *CityServiceImpl) GetCityByID(ctx context.Context, cityID uuid.UUID) (*dto.CityResponse, error) {
	city, err := s.cityRepo.GetByID(ctx, cityID)
	if err != nil {
		log.Printf("Error retrieving city: %v", err)
		return nil, fmt.Errorf("failed to get city: %w", err)
	}

	if city == nil {
		return nil, apperrors.ErrCityNotFound
	}

	return s.toCityResponse(ctx, city, nil), nil
}

func (s *CityServiceImpl) ListCities(ctx context.Context, req dto.CityListRequest) (*dto.CityListResponse, error) {
	var cities []*models.City
	var total int64
	var err error

	// Handle search if provided
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		cities, err = s.cityRepo.SearchByTitle(ctx, strings.TrimSpace(*req.Search), req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error searching cities: %v", err)
			return nil, fmt.Errorf("failed to search cities: %w", err)
		}
		// For search, approximate total
		total = int64(len(cities))
	} else if req.Status != nil && *req.Status == "publish" {
		// Get published cities
		cities, err = s.cityRepo.GetPublished(ctx, req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error getting published cities: %v", err)
			return nil, fmt.Errorf("failed to get published cities: %w", err)
		}
		// Get total count (simplified - could be optimized with a separate count method)
		allPublished, countErr := s.cityRepo.GetPublished(ctx, 0, 1000)
		if countErr != nil {
			total = int64(len(cities))
		} else {
			total = int64(len(allPublished))
		}
	} else {
		// List all cities
		cities, err = s.cityRepo.List(ctx, req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error listing cities: %v", err)
			return nil, fmt.Errorf("failed to list cities: %w", err)
		}

		// Get total count
		total, err = s.cityRepo.Count(ctx)
		if err != nil {
			log.Printf("Error counting cities: %v", err)
			return nil, fmt.Errorf("failed to count cities: %w", err)
		}

		// Filter by draft status if specified
		if req.Status != nil && *req.Status == "draft" {
			filteredCities := make([]*models.City, 0)
			for _, city := range cities {
				if city.Status == models.CityStatusDraft {
					filteredCities = append(filteredCities, city)
				}
			}
			cities = filteredCities
		}
	}

	// Convert to response DTOs
	responseCities := make([]*dto.CityResponse, len(cities))
	for i, city := range cities {
		responseCities[i] = s.toCityResponse(ctx, city, req.Stat)
	}

	pagination := utils.CalculatePagination(total, req.Offset, req.Limit)

	return &dto.CityListResponse{
		Cities:     responseCities,
		Pagination: pagination,
	}, nil
}

// Status Management
func (s *CityServiceImpl) UpdateStatus(ctx context.Context, req dto.UpdateCityStatusRequest) (*dto.CityResponse, error) {
	// Business validation - validate status
	if req.Status != "publish" && req.Status != "draft" {
		return nil, apperrors.ErrInvalidCityStatus
	}

	// Get existing city
	city, err := s.cityRepo.GetByID(ctx, req.CityID)
	if err != nil {
		log.Printf("Error retrieving city: %v", err)
		return nil, fmt.Errorf("failed to get city: %w", err)
	}

	if city == nil {
		return nil, apperrors.ErrCityNotFound
	}

	// Update status
	status := models.CityStatus(req.Status)
	if err := s.cityRepo.UpdateStatus(ctx, req.CityID, status); err != nil {
		log.Printf("Error updating city status: %v", err)
		return nil, fmt.Errorf("failed to update city status: %w", err)
	}

	city.Status = status
	city.UpdatedAt = time.Now()

	return s.toCityResponse(ctx, city, nil), nil
}

func (s *CityServiceImpl) UpdateCityStatus(ctx context.Context, cityID uuid.UUID, req dto.UpdateCityStatusRequest) (*dto.CityResponse, error) {
	// Override the cityID in the request with the one from URL
	req.CityID = cityID
	return s.UpdateStatus(ctx, req)
}

// Public Operations
func (s *CityServiceImpl) GetPublishedCities(ctx context.Context, offset, limit int) ([]*dto.CityResponse, error) {
	cities, err := s.cityRepo.GetPublished(ctx, offset, limit)
	if err != nil {
		log.Printf("Error getting published cities: %v", err)
		return nil, fmt.Errorf("failed to get published cities: %w", err)
	}

	// Convert to response DTOs
	responseCities := make([]*dto.CityResponse, len(cities))
	for i, city := range cities {
		responseCities[i] = s.toCityResponse(ctx, city, nil)
	}

	return responseCities, nil
}

func (s *CityServiceImpl) GetCityTours(ctx context.Context, cityID uuid.UUID, offset, limit int) ([]*dto.TourResponse, error) {
	// Check if city exists
	city, err := s.cityRepo.GetByID(ctx, cityID)
	if err != nil {
		log.Printf("Error retrieving city: %v", err)
		return nil, fmt.Errorf("failed to get city: %w", err)
	}

	if city == nil {
		return nil, apperrors.ErrCityNotFound
	}

	// Get tours for the city
	tours, err := s.tourRepo.GetByCity(ctx, cityID, offset, limit)
	if err != nil {
		log.Printf("Error getting tours for city: %v", err)
		return nil, fmt.Errorf("failed to get city tours: %w", err)
	}

	// Convert to response DTOs
	responseTours := make([]*dto.TourResponse, len(tours))
	for i, tour := range tours {
		responseTours[i] = s.toTourResponse(tour)
	}

	return responseTours, nil
}

// Helper methods
func (s *CityServiceImpl) toCityResponse(ctx context.Context, city *models.City, stat *bool) *dto.CityResponse {
	response := &dto.CityResponse{
		ID:        city.ID,
		Title:     city.Title,
		Detail:    city.Detail,
		Image:     cdn.GetFullFileURL(city.Image), // Construct full URL for frontend
		Status:    string(city.Status),
		CreatedAt: city.CreatedAt.Format(time.RFC3339),
		UpdatedAt: city.UpdatedAt.Format(time.RFC3339),
	}

	// Calculate statistics if requested
	if stat != nil && *stat {
		activeTours, activePeriods := s.calculateCityStatistics(ctx, city.ID)
		response.ActiveTours = &activeTours
		response.ActivePeriods = &activePeriods
	}

	return response
}

// calculateCityStatistics calculates active tours and periods for a city using same logic as ListToursWithPeriods
func (s *CityServiceImpl) calculateCityStatistics(ctx context.Context, cityID uuid.UUID) (int, int) {
	// Get published tours for this city
	tours, err := s.tourRepo.GetByCity(ctx, cityID, 0, 1000) // Get all tours for the city
	if err != nil {
		log.Printf("Error getting tours for city statistics: %v", err)
		return 0, 0
	}

	activeTours := 0
	activePeriods := 0
	now := time.Now()

	for _, tour := range tours {
		// Only count published tours (same logic as tour service)
		if tour.Status != models.TourStatusPublish {
			continue
		}

		// Get periods for this tour
		periods, err := s.periodRepo.GetByTourID(ctx, tour.ID)
		if err != nil {
			log.Printf("Error getting periods for tour %s: %v", tour.ID, err)
			continue
		}

		// Count available periods for this tour (same logic as ListToursWithPeriods)
		tourHasPeriods := false
		for _, period := range periods {
			// Skip draft periods
			if period.Status != models.PeriodStatusPublish {
				continue
			}

			// Skip periods that have already started (same day or past)
			if period.StartDate.Before(now) || period.StartDate.Format("2006-01-02") == now.Format("2006-01-02") {
				continue
			}

			// This is an active/available period
			activePeriods++
			tourHasPeriods = true
		}

		// Only count tour if it has at least one available period
		if tourHasPeriods {
			activeTours++
		}
	}

	return activeTours, activePeriods
}

func (s *CityServiceImpl) toTourResponse(tour *models.Tour) *dto.TourResponse {
	var cityTitle *string
	// Note: This implementation doesn't include city title lookup for performance
	// If needed, this could be enhanced to lookup city title from cityRepo

	return &dto.TourResponse{
		ID:        tour.ID,
		Title:     tour.Title,
		Code:      tour.Code,
		CityID:    tour.CityID,
		CityTitle: cityTitle,
		Detail:    tour.Detail,
		Thumbnail: tour.Thumbnail,
		FilePDF:   tour.FilePDF,
		FileWord:  tour.FileWord,
		Status:    string(tour.Status),
		CreatedAt: tour.CreatedAt.Format(time.RFC3339),
		UpdatedAt: tour.UpdatedAt.Format(time.RFC3339),
	}
}