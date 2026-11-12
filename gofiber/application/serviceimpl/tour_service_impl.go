package serviceimpl

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/cdn"
	apperrors "gofiber-template/pkg/errors"
	"gofiber-template/pkg/utils"

	"github.com/google/uuid"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
)

type TourServiceImpl struct {
	tourRepo      repositories.TourRepository
	cityRepo      repositories.CityRepository
	periodRepo    repositories.PeriodRepository
	bookingRepo   repositories.BookingRepository
	seatHoldRepo  repositories.SeatHoldRepository
}

func NewTourService(
	tourRepo repositories.TourRepository,
	cityRepo repositories.CityRepository,
	periodRepo repositories.PeriodRepository,
	bookingRepo repositories.BookingRepository,
	seatHoldRepo repositories.SeatHoldRepository,
) services.TourService {
	return &TourServiceImpl{
		tourRepo:     tourRepo,
		cityRepo:     cityRepo,
		periodRepo:   periodRepo,
		bookingRepo:  bookingRepo,
		seatHoldRepo: seatHoldRepo,
	}
}

func (s *TourServiceImpl) CreateTour(ctx context.Context, req dto.CreateTourRequest) (*dto.TourResponse, error) {
	// Business validation - validate tour code format
	if !isValidTourCode(req.Code) {
		return nil, apperrors.ErrInvalidTourCode
	}

	// Business validation - check if tour code already exists
	if existing, err := s.tourRepo.GetByCode(ctx, req.Code); err != nil {
		log.Printf("Error checking tour code existence: %v", err)
		return nil, fmt.Errorf("failed to validate tour code: %w", err)
	} else if existing != nil {
		return nil, apperrors.ErrTourCodeExists
	}

	// Business validation - validate city if provided
	if req.CityID != nil {
		if city, err := s.cityRepo.GetByID(ctx, *req.CityID); err != nil {
			log.Printf("Error checking city existence: %v", err)
			return nil, fmt.Errorf("failed to validate city: %w", err)
		} else if city == nil {
			return nil, apperrors.ErrCityNotFound
		}
	}

	// Determine status - use request status or default to draft
	status := models.TourStatusDraft // Default
	if req.Status != nil {
		switch *req.Status {
		case "publish":
			status = models.TourStatusPublish
		case "draft":
			status = models.TourStatusDraft
		default:
			return nil, apperrors.ErrInvalidTourStatus
		}
	}

	// Create domain model
	tour := &models.Tour{
		Title:     req.Title,
		Code:      strings.ToUpper(req.Code),
		CityID:    req.CityID,
		Detail:    req.Detail,
		Thumbnail: req.Thumbnail,
		FilePDF:   req.FilePDF,
		FileWord:  req.FileWord,
		Status:    status, // Use requested status or default
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create tour
	if err := s.tourRepo.Create(ctx, tour); err != nil {
		log.Printf("Error creating tour: %v", err)
		return nil, fmt.Errorf("failed to create tour: %w", err)
	}

	return s.toTourResponse(ctx, tour), nil
}

func (s *TourServiceImpl) GetTourByID(ctx context.Context, tourID uuid.UUID) (*dto.TourResponse, error) {
	tour, err := s.tourRepo.GetByID(ctx, tourID)
	if err != nil {
		log.Printf("Error retrieving tour: %v", err)
		return nil, fmt.Errorf("failed to get tour: %w", err)
	}

	if tour == nil {
		return nil, apperrors.ErrTourNotFound
	}

	return s.toTourResponse(ctx, tour), nil
}

func (s *TourServiceImpl) GetTourByCode(ctx context.Context, code string) (*dto.TourResponse, error) {
	tour, err := s.tourRepo.GetByCode(ctx, strings.ToUpper(code))
	if err != nil {
		log.Printf("Error retrieving tour by code: %v", err)
		return nil, fmt.Errorf("failed to get tour by code: %w", err)
	}

	if tour == nil {
		return nil, apperrors.ErrTourNotFound
	}

	return s.toTourResponse(ctx, tour), nil
}

func (s *TourServiceImpl) GetTourWithPeriods(ctx context.Context, tourID uuid.UUID) (*dto.TourWithPeriodsResponse, error) {
	tour, err := s.tourRepo.GetWithPeriods(ctx, tourID)
	if err != nil {
		log.Printf("Error retrieving tour with periods: %v", err)
		return nil, fmt.Errorf("failed to get tour with periods: %w", err)
	}

	if tour == nil {
		return nil, apperrors.ErrTourNotFound
	}

	return s.toTourWithPeriodsResponse(ctx, tour), nil
}

func (s *TourServiceImpl) UpdateTour(ctx context.Context, tourID uuid.UUID, req dto.UpdateTourRequest) (*dto.TourResponse, error) {
	// Get existing tour
	tour, err := s.tourRepo.GetByID(ctx, tourID)
	if err != nil {
		log.Printf("Error retrieving tour: %v", err)
		return nil, fmt.Errorf("failed to get tour: %w", err)
	}

	if tour == nil {
		return nil, apperrors.ErrTourNotFound
	}

	// Business validation - validate city if provided
	if req.CityID != nil {
		if city, err := s.cityRepo.GetByID(ctx, *req.CityID); err != nil {
			log.Printf("Error checking city existence: %v", err)
			return nil, fmt.Errorf("failed to validate city: %w", err)
		} else if city == nil {
			return nil, apperrors.ErrCityNotFound
		}
	}

	// Business validation - validate status if provided
	if req.Status != nil {
		status := models.TourStatus(*req.Status)
		if status != models.TourStatusPublish && status != models.TourStatusDraft {
			return nil, apperrors.ErrInvalidTourStatus
		}
	}

	// Update fields if provided
	if req.Title != nil {
		tour.Title = *req.Title
	}
	if req.CityID != nil {
		tour.CityID = req.CityID
	}
	if req.Detail != nil {
		tour.Detail = req.Detail
	}
	if req.Thumbnail != nil {
		tour.Thumbnail = req.Thumbnail
	}
	if req.FilePDF != nil {
		tour.FilePDF = req.FilePDF
	}
	if req.FileWord != nil {
		tour.FileWord = req.FileWord
	}
	if req.Status != nil {
		tour.Status = models.TourStatus(*req.Status)
	}

	tour.UpdatedAt = time.Now()

	// Update tour
	if err := s.tourRepo.Update(ctx, tourID, tour); err != nil {
		log.Printf("Error updating tour: %v", err)
		return nil, fmt.Errorf("failed to update tour: %w", err)
	}

	return s.toTourResponse(ctx, tour), nil
}

func (s *TourServiceImpl) DeleteTour(ctx context.Context, tourID uuid.UUID) error {
	// Check if tour exists
	tour, err := s.tourRepo.GetByID(ctx, tourID)
	if err != nil {
		log.Printf("Error retrieving tour: %v", err)
		return fmt.Errorf("failed to get tour: %w", err)
	}

	if tour == nil {
		return apperrors.ErrTourNotFound
	}

	// Business rule - check if tour has active bookings (simplified check)
	// In a real implementation, you would check for active bookings through booking repository

	// Delete tour
	if err := s.tourRepo.Delete(ctx, tourID); err != nil {
		log.Printf("Error deleting tour: %v", err)
		return fmt.Errorf("failed to delete tour: %w", err)
	}

	return nil
}

func (s *TourServiceImpl) ListTours(ctx context.Context, req dto.TourListRequest) (*dto.TourListResponse, error) {
	var tours []*models.Tour
	var total int64
	var err error

	// Handle different filtering options
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		// Search by title or code
		tours, err = s.tourRepo.SearchByTitleOrCode(ctx, strings.TrimSpace(*req.Search), req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error searching tours: %v", err)
			return nil, fmt.Errorf("failed to search tours: %w", err)
		}
		total = int64(len(tours)) // Simplified count for search
	} else if req.CityID != nil {
		// Filter by city
		tours, err = s.tourRepo.GetByCity(ctx, *req.CityID, req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error getting tours by city: %v", err)
			return nil, fmt.Errorf("failed to get tours by city: %w", err)
		}
		total = int64(len(tours)) // Simplified count
	} else if req.Status != nil && *req.Status == "publish" {
		// Get published tours
		tours, err = s.tourRepo.GetPublished(ctx, req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error getting published tours: %v", err)
			return nil, fmt.Errorf("failed to get published tours: %w", err)
		}
		total = int64(len(tours)) // Simplified count
	} else {
		// List all tours
		tours, err = s.tourRepo.List(ctx, req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error listing tours: %v", err)
			return nil, fmt.Errorf("failed to list tours: %w", err)
		}

		// Get total count
		total, err = s.tourRepo.Count(ctx)
		if err != nil {
			log.Printf("Error counting tours: %v", err)
			return nil, fmt.Errorf("failed to count tours: %w", err)
		}
	}

	// Convert to response DTOs
	responseTours := make([]*dto.TourResponse, len(tours))
	for i, tour := range tours {
		responseTours[i] = s.toTourResponse(ctx, tour)
	}

	pagination := dto.PaginationMeta{
		Total:       total,
		HasMore:     int64(req.Offset+len(tours)) < total,
		CurrentPage: (req.Offset / req.Limit) + 1,
		TotalPages:  int((total + int64(req.Limit) - 1) / int64(req.Limit)),
		PageSize:    req.Limit,
		Offset:      req.Offset,
	}

	return &dto.TourListResponse{
		Tours:      responseTours,
		Pagination: pagination,
	}, nil
}

func (s *TourServiceImpl) ListToursWithPeriods(ctx context.Context, req dto.TourListRequest, availableOnly bool) (*dto.TourWithPeriodsListResponse, error) {
	var tours []*models.Tour
	var err error

	// Handle different filtering options (same as ListTours)
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		tours, err = s.tourRepo.SearchByTitleOrCode(ctx, strings.TrimSpace(*req.Search), req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error searching tours with periods: %v", err)
			return nil, fmt.Errorf("failed to search tours with periods: %w", err)
		}
	} else if req.CityID != nil {
		tours, err = s.tourRepo.GetByCity(ctx, *req.CityID, req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error getting tours by city with periods: %v", err)
			return nil, fmt.Errorf("failed to get tours by city with periods: %w", err)
		}
	} else if req.Status != nil && *req.Status == "publish" {
		tours, err = s.tourRepo.GetPublished(ctx, req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error getting published tours with periods: %v", err)
			return nil, fmt.Errorf("failed to get published tours with periods: %w", err)
		}
	} else {
		tours, err = s.tourRepo.List(ctx, req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error listing tours with periods: %v", err)
			return nil, fmt.Errorf("failed to list tours with periods: %w", err)
		}
	}

	// Convert to response DTOs with periods
	responseTours := make([]*dto.TourWithPeriodsResponse, len(tours))
	for i, tour := range tours {
		// Get periods for this tour with preloaded related data
		periods, err := s.periodRepo.GetByTourIDWithPagination(ctx, tour.ID, 0, 100) // Get all periods
		if err != nil {
			log.Printf("Error getting periods for tour %s: %v", tour.ID, err)
			// Continue with empty periods instead of failing
			periods = []*models.Period{}
		}

		// Filter and convert periods
		var availablePeriods []*dto.PeriodResponse
		var minStartDate *string
		now := time.Now()

		for _, period := range periods {
			// Skip draft periods for public routes
			if period.Status != models.PeriodStatusPublish {
				continue
			}

			// Skip periods that have already started (same day or past)
			if period.StartDate.Before(now) || period.StartDate.Format("2006-01-02") == now.Format("2006-01-02") {
				continue
			}

			// Convert using enhanced toPeriodResponse that handles related tables properly
			periodResp := s.toEnhancedPeriodResponse(ctx, period)

			// If availableOnly is true, skip fully booked periods (use calculated available seats)
			if availableOnly && periodResp.AvailableSeats <= 0 {
				continue
			}

			availablePeriods = append(availablePeriods, periodResp)

			// Track minimum start date
			if minStartDate == nil || periodResp.StartDate < *minStartDate {
				minStartDate = &periodResp.StartDate
			}
		}

		// Only include tours that have available periods
		if len(availablePeriods) > 0 {
			responseTours[i] = &dto.TourWithPeriodsResponse{
				TourResponse:     s.toTourResponse(ctx, tour),
				AvailablePeriods: availablePeriods,
				MinStartDate:     minStartDate,
			}
		} else {
			// Mark as nil to filter out later
			responseTours[i] = nil
		}
	}

	// Filter out nil tours (tours without available periods)
	var filteredTours []*dto.TourWithPeriodsResponse
	for _, tour := range responseTours {
		if tour != nil {
			filteredTours = append(filteredTours, tour)
		}
	}

	// Update total count to reflect actual available tours
	actualTotal := int64(len(filteredTours))

	pagination := dto.PaginationMeta{
		Total:       actualTotal,
		HasMore:     false, // Since we're filtering, disable pagination temporarily
		CurrentPage: 1,
		TotalPages:  1,
		PageSize:    req.Limit,
		Offset:      0,
	}

	return &dto.TourWithPeriodsListResponse{
		Tours:      filteredTours,
		Pagination: pagination,
	}, nil
}

func (s *TourServiceImpl) ListSpecialTours(ctx context.Context, req dto.TourListRequest, availableOnly bool) (*dto.TourWithPeriodsListResponse, error) {
	var tours []*models.Tour
	var err error

	// ดึงเฉพาะทัวร์ที่มี sale periods จาก DB ตรงๆ (ไม่ filter ใน application layer)
	tours, err = s.tourRepo.GetToursWithSalePeriods(ctx, req.Offset, req.Limit)
	if err != nil {
		log.Printf("Error getting special tours: %v", err)
		return nil, fmt.Errorf("failed to get special tours: %w", err)
	}

	// Convert to response DTOs with periods (only special periods)
	responseTours := make([]*dto.TourWithPeriodsResponse, len(tours))
	for i, tour := range tours {
		// Get periods for this tour with preloaded related data
		periods, err := s.periodRepo.GetByTourIDWithPagination(ctx, tour.ID, 0, 100) // Get all periods
		if err != nil {
			log.Printf("Error getting periods for special tour %s: %v", tour.ID, err)
			// Continue with empty periods instead of failing
			periods = []*models.Period{}
		}

		// Filter and convert periods - only include special periods (sale rate/commission)
		var specialPeriods []*dto.PeriodResponse
		var minStartDate *string
		now := time.Now()

		for _, period := range periods {
			// Skip draft periods for public routes
			if period.Status != models.PeriodStatusPublish {
				continue
			}

			// Skip periods that have already started (same day or past)
			if period.StartDate.Before(now) || period.StartDate.Format("2006-01-02") == now.Format("2006-01-02") {
				continue
			}

			// Check if this period has sale rate type or sale commission type
			hasSaleRate := string(period.RateType) == "sale"
			hasSaleCommission := string(period.CommissionType) == "sale"

			// Skip periods that don't have sale rate or commission
			if !hasSaleRate && !hasSaleCommission {
				continue
			}

			// Convert using enhanced toPeriodResponse that handles related tables properly
			periodResp := s.toEnhancedPeriodResponse(ctx, period)

			// If availableOnly is true, skip fully booked periods (use calculated available seats)
			if availableOnly && periodResp.AvailableSeats <= 0 {
				continue
			}

			specialPeriods = append(specialPeriods, periodResp)

			// Track minimum start date
			if minStartDate == nil || periodResp.StartDate < *minStartDate {
				minStartDate = &periodResp.StartDate
			}
		}

		// Only include tours that have special periods
		if len(specialPeriods) > 0 {
			responseTours[i] = &dto.TourWithPeriodsResponse{
				TourResponse:     s.toTourResponse(ctx, tour),
				AvailablePeriods: specialPeriods,
				MinStartDate:     minStartDate,
			}
		} else {
			// Mark as nil to filter out later
			responseTours[i] = nil
		}
	}

	// Filter out nil tours (tours without special periods)
	var filteredTours []*dto.TourWithPeriodsResponse
	for _, tour := range responseTours {
		if tour != nil {
			filteredTours = append(filteredTours, tour)
		}
	}

	// Update total count to reflect actual available special tours
	actualTotal := int64(len(filteredTours))

	pagination := dto.PaginationMeta{
		Total:       actualTotal,
		HasMore:     false, // Since we're filtering, disable pagination temporarily
		CurrentPage: 1,
		TotalPages:  1,
		PageSize:    req.Limit,
		Offset:      0,
	}

	return &dto.TourWithPeriodsListResponse{
		Tours:      filteredTours,
		Pagination: pagination,
	}, nil
}

func (s *TourServiceImpl) GetAvailableTours(ctx context.Context, offset, limit int) ([]*dto.TourWithPeriodsResponse, error) {
	tours, err := s.tourRepo.GetPublishedWithAvailablePeriods(ctx, offset, limit)
	if err != nil {
		log.Printf("Error getting available tours: %v", err)
		return nil, fmt.Errorf("failed to get available tours: %w", err)
	}

	// Convert to response DTOs
	responseTours := make([]*dto.TourWithPeriodsResponse, len(tours))
	for i, tour := range tours {
		responseTours[i] = s.toTourWithPeriodsResponse(ctx, tour)
	}

	return responseTours, nil
}

func (s *TourServiceImpl) GetToursByCity(ctx context.Context, cityID uuid.UUID, offset, limit int) ([]*dto.TourResponse, error) {
	tours, err := s.tourRepo.GetByCity(ctx, cityID, offset, limit)
	if err != nil {
		log.Printf("Error getting tours by city: %v", err)
		return nil, fmt.Errorf("failed to get tours by city: %w", err)
	}

	// Convert to response DTOs
	responseTours := make([]*dto.TourResponse, len(tours))
	for i, tour := range tours {
		responseTours[i] = s.toTourResponse(ctx, tour)
	}

	return responseTours, nil
}

func (s *TourServiceImpl) GetPublishedTours(ctx context.Context, offset, limit int) ([]*dto.TourResponse, error) {
	tours, err := s.tourRepo.GetPublished(ctx, offset, limit)
	if err != nil {
		log.Printf("Error getting published tours: %v", err)
		return nil, fmt.Errorf("failed to get published tours: %w", err)
	}

	// Convert to response DTOs
	responseTours := make([]*dto.TourResponse, len(tours))
	for i, tour := range tours {
		responseTours[i] = s.toTourResponse(ctx, tour)
	}

	return responseTours, nil
}

func (s *TourServiceImpl) SearchTours(ctx context.Context, query string, offset, limit int) ([]*dto.TourResponse, error) {
	// Business validation - clean search query
	query = strings.TrimSpace(query)
	if query == "" {
		return []*dto.TourResponse{}, nil
	}

	tours, err := s.tourRepo.SearchByTitle(ctx, query, offset, limit)
	if err != nil {
		log.Printf("Error searching tours: %v", err)
		return nil, fmt.Errorf("failed to search tours: %w", err)
	}

	// Convert to response DTOs
	responseTours := make([]*dto.TourResponse, len(tours))
	for i, tour := range tours {
		responseTours[i] = s.toTourResponse(ctx, tour)
	}

	return responseTours, nil
}

func (s *TourServiceImpl) UploadTourFile(ctx context.Context, req dto.UploadTourFileRequest, fileType string) (*dto.TourResponse, error) {
	// Business validation - check file size
	if len(req.FileData) > MaxFileSize {
		return nil, apperrors.ErrTourFileTooLarge
	}

	// Business validation - check file type
	if !isValidTourFileType(req.FileType, fileType) {
		return nil, apperrors.ErrInvalidTourFileType
	}

	// Get existing tour
	tour, err := s.tourRepo.GetByID(ctx, req.TourID)
	if err != nil {
		log.Printf("Error retrieving tour: %v", err)
		return nil, fmt.Errorf("failed to get tour: %w", err)
	}

	if tour == nil {
		return nil, apperrors.ErrTourNotFound
	}

	// Generate file URL/path (in real implementation, upload to storage service)
	fileURL := s.generateFileURL(req.TourID, req.FileName, fileType)

	// Update appropriate file field
	switch fileType {
	case "pdf":
		tour.FilePDF = &fileURL
	case "word":
		tour.FileWord = &fileURL
	default:
		return nil, apperrors.ErrInvalidTourFileType
	}

	tour.UpdatedAt = time.Now()

	// Update tour
	if err := s.tourRepo.Update(ctx, req.TourID, tour); err != nil {
		log.Printf("Error updating tour file: %v", err)
		return nil, fmt.Errorf("failed to update tour file: %w", err)
	}

	return s.toTourResponse(ctx, tour), nil
}

func (s *TourServiceImpl) DeleteTourFile(ctx context.Context, tourID uuid.UUID, fileType string) (*dto.TourResponse, error) {
	// Get existing tour
	tour, err := s.tourRepo.GetByID(ctx, tourID)
	if err != nil {
		log.Printf("Error retrieving tour: %v", err)
		return nil, fmt.Errorf("failed to get tour: %w", err)
	}

	if tour == nil {
		return nil, apperrors.ErrTourNotFound
	}

	// Remove appropriate file field
	switch fileType {
	case "pdf":
		tour.FilePDF = nil
	case "word":
		tour.FileWord = nil
	default:
		return nil, apperrors.ErrInvalidTourFileType
	}

	tour.UpdatedAt = time.Now()

	// Update tour
	if err := s.tourRepo.Update(ctx, tourID, tour); err != nil {
		log.Printf("Error removing tour file: %v", err)
		return nil, fmt.Errorf("failed to remove tour file: %w", err)
	}

	return s.toTourResponse(ctx, tour), nil
}

func (s *TourServiceImpl) PublishTour(ctx context.Context, tourID uuid.UUID) (*dto.TourResponse, error) {
	// Get existing tour
	tour, err := s.tourRepo.GetByID(ctx, tourID)
	if err != nil {
		log.Printf("Error retrieving tour: %v", err)
		return nil, fmt.Errorf("failed to get tour: %w", err)
	}

	if tour == nil {
		return nil, apperrors.ErrTourNotFound
	}

	// Update status to published
	if err := s.tourRepo.UpdateStatus(ctx, tourID, models.TourStatusPublish); err != nil {
		log.Printf("Error publishing tour: %v", err)
		return nil, fmt.Errorf("failed to publish tour: %w", err)
	}

	tour.Status = models.TourStatusPublish
	tour.UpdatedAt = time.Now()

	return s.toTourResponse(ctx, tour), nil
}

func (s *TourServiceImpl) UnpublishTour(ctx context.Context, tourID uuid.UUID) (*dto.TourResponse, error) {
	// Get existing tour
	tour, err := s.tourRepo.GetByID(ctx, tourID)
	if err != nil {
		log.Printf("Error retrieving tour: %v", err)
		return nil, fmt.Errorf("failed to get tour: %w", err)
	}

	if tour == nil {
		return nil, apperrors.ErrTourNotFound
	}

	// Update status to draft
	if err := s.tourRepo.UpdateStatus(ctx, tourID, models.TourStatusDraft); err != nil {
		log.Printf("Error unpublishing tour: %v", err)
		return nil, fmt.Errorf("failed to unpublish tour: %w", err)
	}

	tour.Status = models.TourStatusDraft
	tour.UpdatedAt = time.Now()

	return s.toTourResponse(ctx, tour), nil
}

func (s *TourServiceImpl) GetTourStats(ctx context.Context) (*dto.TourStatsResponse, error) {
	// Get total tours count
	totalTours, err := s.tourRepo.Count(ctx)
	if err != nil {
		log.Printf("Error getting total tours count: %v", err)
		return nil, fmt.Errorf("failed to get tour stats: %w", err)
	}

	// Get published tours (simplified - in real implementation use dedicated methods)
	publishedTours, err := s.tourRepo.GetPublished(ctx, 0, 1000)
	if err != nil {
		log.Printf("Error getting published tours for stats: %v", err)
		return nil, fmt.Errorf("failed to get tour stats: %w", err)
	}

	publishedCount := int64(len(publishedTours))
	draftCount := totalTours - publishedCount

	// Tours with periods (simplified calculation)
	toursWithPeriods := int64(0) // In real implementation, join with periods table

	return &dto.TourStatsResponse{
		TotalTours:       totalTours,
		PublishedTours:   publishedCount,
		DraftTours:       draftCount,
		ToursWithPeriods: toursWithPeriods,
	}, nil
}

func (s *TourServiceImpl) ValidateTourCode(ctx context.Context, code string, excludeTourID *uuid.UUID) (bool, error) {
	// Business validation - validate format first
	if !isValidTourCode(code) {
		return false, nil
	}

	// Check if tour code exists
	existing, err := s.tourRepo.GetByCode(ctx, strings.ToUpper(code))
	if err != nil {
		log.Printf("Error checking tour code existence: %v", err)
		return false, fmt.Errorf("failed to validate tour code: %w", err)
	}

	// If no existing tour found, code is available
	if existing == nil {
		return true, nil
	}

	// If excludeTourID is provided and matches, code is available for that tour
	if excludeTourID != nil && existing.ID == *excludeTourID {
		return true, nil
	}

	// Code already exists for another tour
	return false, nil
}

// UpdateTourStatus updates tour status
func (s *TourServiceImpl) UpdateTourStatus(ctx context.Context, req dto.UpdateTourStatusRequest) (*dto.TourResponse, error) {
	// Get existing tour
	tour, err := s.tourRepo.GetByID(ctx, req.ID)
	if err != nil {
		log.Printf("Error retrieving tour: %v", err)
		return nil, fmt.Errorf("failed to get tour: %w", err)
	}

	if tour == nil {
		return nil, apperrors.ErrTourNotFound
	}

	// Business validation - check if status change is needed
	currentStatus := string(tour.Status)
	if currentStatus == req.Status {
		return nil, fmt.Errorf("tour is already in %s status", req.Status)
	}

	// Validate status value
	var newStatus models.TourStatus
	switch req.Status {
	case "publish":
		newStatus = models.TourStatusPublish
	case "draft":
		newStatus = models.TourStatusDraft
	default:
		return nil, fmt.Errorf("invalid tour status: %s", req.Status)
	}

	// Update status
	if err := s.tourRepo.UpdateStatus(ctx, req.ID, newStatus); err != nil {
		log.Printf("Error updating tour status: %v", err)
		return nil, fmt.Errorf("failed to update tour status: %w", err)
	}

	tour.Status = newStatus
	tour.UpdatedAt = time.Now()

	return s.toTourResponse(ctx, tour), nil
}

// GetTourPeriods returns all periods for a tour
func (s *TourServiceImpl) GetTourPeriods(ctx context.Context, tourID uuid.UUID) ([]*dto.PeriodResponse, error) {
	// Validate tour exists
	tour, err := s.tourRepo.GetByID(ctx, tourID)
	if err != nil {
		log.Printf("Error retrieving tour: %v", err)
		return nil, fmt.Errorf("failed to get tour: %w", err)
	}

	if tour == nil {
		return nil, apperrors.ErrTourNotFound
	}

	// Get periods for tour
	periods, err := s.periodRepo.GetByTourID(ctx, tourID) // Get all periods
	if err != nil {
		log.Printf("Error getting tour periods: %v", err)
		return nil, fmt.Errorf("failed to get tour periods: %w", err)
	}

	// Convert to DTO responses
	periodResponses := make([]*dto.PeriodResponse, len(periods))
	for i, period := range periods {
		periodResponses[i] = s.toPeriodResponse(ctx, period)
	}

	return periodResponses, nil
}

func (s *TourServiceImpl) GetTourPeriodsWithPagination(ctx context.Context, tourID uuid.UUID, offset, limit int) (*dto.PeriodListResponse, error) {
	// Validate tour exists
	tour, err := s.tourRepo.GetByID(ctx, tourID)
	if err != nil {
		log.Printf("Error retrieving tour: %v", err)
		return nil, fmt.Errorf("failed to get tour: %w", err)
	}

	if tour == nil {
		return nil, apperrors.ErrTourNotFound
	}

	// Get periods for tour with pagination
	periods, err := s.periodRepo.GetByTourIDWithPagination(ctx, tourID, offset, limit)
	if err != nil {
		log.Printf("Error getting tour periods: %v", err)
		return nil, fmt.Errorf("failed to get tour periods: %w", err)
	}

	// Get total count for pagination
	totalCount, err := s.periodRepo.CountByTourID(ctx, tourID)
	if err != nil {
		log.Printf("Error counting tour periods: %v", err)
		return nil, fmt.Errorf("failed to count tour periods: %w", err)
	}

	// Convert to DTO responses
	periodResponses := make([]*dto.PeriodResponse, len(periods))
	for i, period := range periods {
		periodResponses[i] = s.toPeriodResponse(ctx, period)
	}

	return &dto.PeriodListResponse{
		Periods:    periodResponses,
		Pagination: utils.CalculatePagination(totalCount, offset, limit),
	}, nil
}

func (s *TourServiceImpl) toTourResponse(ctx context.Context, tour *models.Tour) *dto.TourResponse {
	response := &dto.TourResponse{
		ID:        tour.ID,
		Title:     tour.Title,
		Code:      tour.Code,
		CityID:    tour.CityID,
		Detail:    tour.Detail,
		Thumbnail: cdn.GetFullFileURLPtr(tour.Thumbnail), // Construct full URL for frontend
		FilePDF:   cdn.GetFullFileURLPtr(tour.FilePDF),   // Construct full URL for frontend
		FileWord:  cdn.GetFullFileURLPtr(tour.FileWord),  // Construct full URL for frontend
		Status:    string(tour.Status),
		CreatedAt: tour.CreatedAt.Format(time.RFC3339),
		UpdatedAt: tour.UpdatedAt.Format(time.RFC3339),
	}

	// Get city title if city ID is provided
	if tour.CityID != nil {
		if city, err := s.cityRepo.GetByID(ctx, *tour.CityID); err == nil && city != nil {
			response.CityTitle = &city.Title
		}
	}

	return response
}

// toPeriodResponse converts a Period model to PeriodResponse DTO with held seats calculation
func (s *TourServiceImpl) toPeriodResponse(ctx context.Context, period *models.Period) *dto.PeriodResponse {
	// Calculate held seats from active seat holds
	heldSeats := s.calculateHeldSeats(ctx, period.ID)

	// Calculate real-time booking count (active bookings: pending, processing, confirm)
	activeBookings := s.calculateBookedPassengers(ctx, period.ID)

	// Calculate real available seats (total seats - active bookings - held)
	availableSeats := period.Seat - activeBookings - heldSeats
	if availableSeats < 0 {
		availableSeats = 0
	}

	response := &dto.PeriodResponse{
		ID:             period.ID,
		TourID:         period.TourID,
		StartDate:      period.StartDate.Format("2006-01-02"),
		EndDate:        period.EndDate.Format("2006-01-02"),
		Seat:           period.Seat,
		Booking:        activeBookings,
		HeldSeats:      heldSeats,
		AvailableSeats: availableSeats,
		Deposit:        period.Deposit,
		DepositType:    period.DepositType,
		Airline:        period.Airline,
		Airport:        period.Airport,
		Note:           period.Note,
		RateType:       string(period.RateType),
		CommissionType:   string(period.CommissionType),
		Status:           string(period.Status),
		ServiceFeePerPax: period.ServiceFeePerPax,
		CreatedAt:        period.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        period.UpdatedAt.Format(time.RFC3339),
		// Expiration Override
		ExpirationOverride: period.ExpirationOverride,
		ExpirationMinutes:  period.ExpirationMinutes,
		ExpirationHours:    period.ExpirationHours,
		ExpirationDays:     period.ExpirationDays,
		ExpirationMessage:  period.ExpirationMessage,
		ExpirationDisplay:  period.GetExpirationDisplay(),
	}

	// Convert JSONB arrays to response DTOs
	response.Sector1 = s.toFlightSectorResponses(period.Sector1)
	response.Sector2 = s.toFlightSectorResponses(period.Sector2)
	response.RateItems = s.toRateItemResponses(period.RateItems)
	response.CommissionItems = s.toCommissionItemResponses(period.CommissionItems)
	response.RoomTypes = s.toRoomTypeResponses(period.RoomTypes)
	response.AdditionalServices = s.toAdditionalServiceResponses(period.AdditionalServices)
	response.Files = s.toFileAttachmentResponses(period.Files)

	// Check WL: มี booking active ที่จ่ายเงินไม่ครบหรือไม่
	response.HasWaitingList = s.checkHasWaitingList(ctx, period.ID)

	return response
}

// toEnhancedPeriodResponse converts Period to PeriodResponse with proper related table handling and real-time calculations
func (s *TourServiceImpl) toEnhancedPeriodResponse(ctx context.Context, period *models.Period) *dto.PeriodResponse {
	// Calculate real-time booking stats using same logic as PeriodService
	realTimeBookedSeats := s.calculateBookedPassengers(ctx, period.ID)
	heldSeats := s.calculateHeldSeats(ctx, period.ID)

	// Calculate real available seats using same formula as PeriodService
	availableSeats := period.Seat - realTimeBookedSeats - heldSeats
	if availableSeats < 0 {
		availableSeats = 0
	}

	response := &dto.PeriodResponse{
		ID:             period.ID,
		TourID:         period.TourID,
		StartDate:      period.StartDate.Format("2006-01-02"),
		EndDate:        period.EndDate.Format("2006-01-02"),
		Seat:           period.Seat,
		Booking:        realTimeBookedSeats,
		HeldSeats:      heldSeats,
		AvailableSeats: availableSeats,
		Deposit:        period.Deposit,
		DepositType:    period.DepositType,
		Airline:        period.Airline,
		Airport:        period.Airport,
		Note:           period.Note,
		RateType:       string(period.RateType),
		CommissionType:   string(period.CommissionType),
		Status:           string(period.Status),
		ServiceFeePerPax: period.ServiceFeePerPax,
		CreatedAt:        period.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        period.UpdatedAt.Format(time.RFC3339),
		// Expiration Override
		ExpirationOverride: period.ExpirationOverride,
		ExpirationMinutes:  period.ExpirationMinutes,
		ExpirationHours:    period.ExpirationHours,
		ExpirationDays:     period.ExpirationDays,
		ExpirationMessage:  period.ExpirationMessage,
		ExpirationDisplay:  period.GetExpirationDisplay(),
	}

	// Handle related tables (PeriodRateItems, PeriodRoomTypes, etc.) with fallback to JSONB
	response.Sector1 = s.toFlightSectorResponses(period.Sector1)
	response.Sector2 = s.toFlightSectorResponses(period.Sector2)
	response.CommissionItems = s.toCommissionItemResponses(period.CommissionItems)
	response.Files = s.toFileAttachmentResponses(period.Files)

	// Rate items: use related tables if available, fallback to JSONB
	if len(period.PeriodRateItems) > 0 {
		response.RateItems = s.toPeriodRateItemResponses(period.PeriodRateItems)
	} else {
		response.RateItems = s.toRateItemResponses(period.RateItems)
	}

	// Room types: use related tables if available, fallback to JSONB
	if len(period.PeriodRoomTypes) > 0 {
		response.RoomTypes = s.toPeriodRoomTypeResponses(period.PeriodRoomTypes)
	} else {
		response.RoomTypes = s.toRoomTypeResponses(period.RoomTypes)
	}

	// Additional services: use related tables if available, fallback to JSONB
	if len(period.PeriodAdditionalServices) > 0 {
		response.AdditionalServices = s.toPeriodAdditionalServiceResponses(period.PeriodAdditionalServices)
	} else {
		response.AdditionalServices = s.toAdditionalServiceResponses(period.AdditionalServices)
	}

	// Check WL
	response.HasWaitingList = s.checkHasWaitingList(ctx, period.ID)

	return response
}

// checkHasWaitingList เช็คว่า period มี booking active ที่จ่ายเงินไม่ครบ
func (s *TourServiceImpl) checkHasWaitingList(ctx context.Context, periodID uuid.UUID) bool {
	if s.bookingRepo == nil {
		return false
	}

	bookings, err := s.bookingRepo.GetByPeriodID(ctx, periodID, 0, 1000)
	if err != nil {
		return false
	}

	for _, booking := range bookings {
		if booking.Status != "pending" && booking.Status != "processing" && booking.Status != "confirm" {
			continue
		}
		total := booking.GrandTotal
		if total == 0 {
			total = booking.SumNet
		}
		if booking.SumPayment < total {
			return true
		}
	}

	return false
}

// Helper methods for converting related tables
func (s *TourServiceImpl) toPeriodRateItemResponses(items []models.PeriodRateItem) []*dto.RateItemResponse {
	responses := make([]*dto.RateItemResponse, len(items))
	for i, item := range items {
		responses[i] = &dto.RateItemResponse{
			ID:        item.ID.String(),
			Title:     item.Title,
			Regular:   item.Regular,
			Sale:      item.Sale,
			Exclude:   item.Exclude,
			SortOrder: item.SortOrder,
		}
	}
	return responses
}

func (s *TourServiceImpl) toPeriodRoomTypeResponses(roomTypes []models.PeriodRoomType) []*dto.RoomTypeResponse {
	responses := make([]*dto.RoomTypeResponse, len(roomTypes))
	for i, roomType := range roomTypes {
		responses[i] = &dto.RoomTypeResponse{
			ID:        roomType.ID.String(),
			Title:     roomType.Title,
			Price:     roomType.Price,
			SortOrder: roomType.SortOrder,
		}
	}
	return responses
}

func (s *TourServiceImpl) toPeriodAdditionalServiceResponses(services []models.PeriodAdditionalService) []*dto.AdditionalServiceResponse {
	responses := make([]*dto.AdditionalServiceResponse, len(services))
	for i, service := range services {
		responses[i] = &dto.AdditionalServiceResponse{
			ID:        service.ID.String(),
			Title:     service.Title,
			Price:     service.Price,
			RateType:  service.RateType,
			UnitName:  service.UnitName,
			SortOrder: service.SortOrder,
		}
	}
	return responses
}

// calculateBookedPassengers calculates real-time booked passengers using same logic as PeriodService
func (s *TourServiceImpl) calculateBookedPassengers(ctx context.Context, periodID uuid.UUID) int {
	if s.bookingRepo == nil {
		log.Printf("[WARNING] BookingRepository not available - using static booking count")
		return 0 // Fallback to 0 if no booking repo
	}

	// Use same active statuses as PeriodService
	activeStatuses := []string{"pending", "processing", "confirm"}
	totalBookedPassengers := 0

	for _, status := range activeStatuses {
		count, passengers, err := s.countBookingsByPeriodAndStatus(ctx, periodID, status)
		if err != nil {
			log.Printf("[WARNING] Failed to count bookings for status %s: %v", status, err)
			continue
		}
		totalBookedPassengers += passengers
		log.Printf("[DEBUG] Status %s: %d bookings, %d passengers", status, count, passengers)
	}

	log.Printf("[DEBUG] Total booked passengers for period %s: %d", periodID.String(), totalBookedPassengers)
	return totalBookedPassengers
}

// calculateHeldSeats calculates active held seats using same logic as PeriodService
func (s *TourServiceImpl) calculateHeldSeats(ctx context.Context, periodID uuid.UUID) int {
	if s.seatHoldRepo == nil {
		log.Printf("[WARNING] SeatHoldRepository not available - held seats = 0")
		return 0 // Fallback to 0 if no seat hold repo
	}

	activeHeldSeats, err := s.seatHoldRepo.GetActiveHeldSeats(ctx, periodID)
	if err != nil {
		log.Printf("[WARNING] Failed to get held seats for period %s: %v", periodID.String(), err)
		return 0 // Fallback to 0 on error
	}

	log.Printf("[DEBUG] Active held seats for period %s: %d", periodID.String(), activeHeldSeats)
	return activeHeldSeats
}

// countBookingsByPeriodAndStatus counts bookings by period and status (same as PeriodService)
func (s *TourServiceImpl) countBookingsByPeriodAndStatus(ctx context.Context, periodID uuid.UUID, status string) (int, int, error) {
	// Use the same repository method as PeriodService
	return s.bookingRepo.CountBookingsByPeriodAndStatus(ctx, periodID, status)
}

// Helper methods to convert JSONB arrays to response DTOs
func (s *TourServiceImpl) toFlightSectorResponses(sectors models.FlightSectors) []*dto.FlightSectorResponse {
	responses := make([]*dto.FlightSectorResponse, len(sectors))
	for i, sector := range sectors {
		responses[i] = &dto.FlightSectorResponse{
			Date:   sector.Date,
			Flight: sector.Flight,
			Route:  sector.Route,
		}
	}
	return responses
}

func (s *TourServiceImpl) toRateItemResponses(items models.RateItems) []*dto.RateItemResponse {
	responses := make([]*dto.RateItemResponse, len(items))
	for i, item := range items {
		responses[i] = &dto.RateItemResponse{
			ID:        item.ID,
			Title:     item.Title,
			Regular:   item.Regular,
			Sale:      item.Sale,
			Exclude:   item.Exclude,
			SortOrder: 0, // Default sort order, could be enhanced
		}
	}
	return responses
}

func (s *TourServiceImpl) toCommissionItemResponses(items models.CommissionItems) []*dto.CommissionItemResponse {
	responses := make([]*dto.CommissionItemResponse, len(items))
	for i, item := range items {
		responses[i] = &dto.CommissionItemResponse{
			Regular: item.Regular,
			Sale:    item.Sale,
		}
	}
	return responses
}

func (s *TourServiceImpl) toRoomTypeResponses(roomTypes models.RoomTypes) []*dto.RoomTypeResponse {
	responses := make([]*dto.RoomTypeResponse, len(roomTypes))
	for i, roomType := range roomTypes {
		responses[i] = &dto.RoomTypeResponse{
			ID:        roomType.ID,
			Title:     roomType.Title,
			Price:     roomType.Price,
			SortOrder: 0, // Default sort order, could be enhanced
		}
	}
	return responses
}

func (s *TourServiceImpl) toAdditionalServiceResponses(services models.AdditionalServices) []*dto.AdditionalServiceResponse {
	responses := make([]*dto.AdditionalServiceResponse, len(services))
	for i, service := range services {
		responses[i] = &dto.AdditionalServiceResponse{
			ID:        service.ID,
			Title:     service.Title,
			Price:     service.Price,
			RateType:  service.RateType,
			UnitName:  service.UnitName,
			SortOrder: 0, // Default sort order, could be enhanced
		}
	}
	return responses
}

func (s *TourServiceImpl) toFileAttachmentResponses(files models.FileAttachments) []*dto.FileAttachmentResponse {
	responses := make([]*dto.FileAttachmentResponse, len(files))
	for i, file := range files {
		responses[i] = &dto.FileAttachmentResponse{
			Name: file.Name,
			Type: file.Type,
			URL:  file.URL,
		}
	}
	return responses
}

func (s *TourServiceImpl) toTourWithPeriodsResponse(ctx context.Context, tour *models.Tour) *dto.TourWithPeriodsResponse {
	response := &dto.TourWithPeriodsResponse{
		TourResponse: s.toTourResponse(ctx, tour),
	}

	// Convert periods to full responses using our toPeriodResponse method
	if len(tour.Periods) > 0 {
		periods := make([]*dto.PeriodResponse, 0)
		var minStartDate *time.Time

		for _, period := range tour.Periods {
			// Only include published periods
			if period.IsPublished() {
				periodResponse := s.toPeriodResponse(ctx, &period)
				periods = append(periods, periodResponse)

				// Track minimum start date
				if minStartDate == nil || period.StartDate.Before(*minStartDate) {
					minStartDate = &period.StartDate
				}
			}
		}

		response.AvailablePeriods = periods
		if minStartDate != nil {
			minDateStr := minStartDate.Format("2006-01-02")
			response.MinStartDate = &minDateStr
		}
	}

	return response
}

func isValidTourCode(code string) bool {
	// Tour code should be 3-20 characters, alphanumeric and hyphens only
	if len(code) < 3 || len(code) > 20 {
		return false
	}

	for _, char := range code {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-') {
			return false
		}
	}

	return true
}

func isValidTourFileType(mimeType, expectedType string) bool {
	switch expectedType {
	case "pdf":
		return mimeType == "application/pdf"
	case "word":
		return mimeType == "application/msword" ||
			mimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return false
	}
}

func (s *TourServiceImpl) generateFileURL(tourID uuid.UUID, fileName, fileType string) string {
	// Generate a unique file path
	ext := filepath.Ext(fileName)
	if ext == "" {
		switch fileType {
		case "pdf":
			ext = ".pdf"
		case "word":
			ext = ".docx"
		}
	}

	timestamp := time.Now().Unix()
	return fmt.Sprintf("tours/%s/%s/%d%s", tourID.String(), fileType, timestamp, ext)
}
