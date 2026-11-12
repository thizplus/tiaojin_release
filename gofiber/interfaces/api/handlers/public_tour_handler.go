package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type PublicTourHandler struct {
	tourService services.TourService
}

func NewPublicTourHandler(tourService services.TourService) *PublicTourHandler {
	return &PublicTourHandler{
		tourService: tourService,
	}
}

// ListAvailableTours returns published tours with available periods (no commission info)
func (h *PublicTourHandler) ListAvailableTours(c *fiber.Ctx) error {
	// Parse query parameters
	pagination := utils.ParsePaginationParams(c, 20)
	search := c.Query("search")
	cityID := c.Query("city_id")

	// Build request
	publishStatus := "publish" // Only show published tours
	req := dto.TourListRequest{
		Status:      &publishStatus,
		WithPeriods: true,
		Offset:      pagination.Offset,
		Limit:       pagination.Limit,
	}

	if search != "" {
		req.Search = &search
	}

	if cityID != "" {
		if cityUUID, err := uuid.Parse(cityID); err == nil {
			req.CityID = &cityUUID
		}
	}

	// Get tours with periods from service (availableOnly = false, we filter by status instead)
	result, err := h.tourService.ListToursWithPeriods(c.Context(), req, false)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve tours", err)
	}

	// Convert to public response (exclude commission info only)
	publicTours := make([]*dto.PublicTourWithPeriodsResponse, 0, len(result.Tours))
	for _, tour := range result.Tours {
		publicTour := convertToPublicTourWithPeriods(tour)
		publicTours = append(publicTours, publicTour)
	}

	publicResponse := &dto.PublicTourListResponse{
		Tours:      publicTours,
		Pagination: result.Pagination,
	}

	return utils.SuccessResponse(c, "Tours retrieved successfully", publicResponse)
}

// GetByID returns a specific tour with periods (no commission info)
func (h *PublicTourHandler) GetByID(c *fiber.Ctx) error {
	tourID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid tour ID")
	}

	// Get tour with periods
	tour, err := h.tourService.GetTourWithPeriods(c.Context(), tourID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Tour not found", err)
	}

	// Only return published tours for public API
	if tour.Status != "publish" {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Tour not found", nil)
	}

	// Convert to public response
	publicTour := convertToPublicTourWithPeriods(tour)

	return utils.SuccessResponse(c, "Tour retrieved successfully", publicTour)
}

// Helper function to convert internal tour response to public response (excludes commission only)
func convertToPublicTourWithPeriods(tour *dto.TourWithPeriodsResponse) *dto.PublicTourWithPeriodsResponse {
	publicTour := &dto.PublicTourWithPeriodsResponse{
		PublicTourResponse: &dto.PublicTourResponse{
			ID:        tour.ID,
			Title:     tour.Title,
			Code:      tour.Code,
			CityID:    tour.CityID,
			CityTitle: tour.CityTitle,
			Detail:    tour.Detail,
			Thumbnail: tour.Thumbnail,
			FilePDF:   tour.FilePDF,
			FileWord:  tour.FileWord,
			Status:    tour.Status,
			CreatedAt: tour.CreatedAt,
			UpdatedAt: tour.UpdatedAt,
		},
		MinStartDate: tour.MinStartDate,
	}

	// Convert periods (exclude commission info only)
	if len(tour.AvailablePeriods) > 0 {
		publicPeriods := make([]*dto.PublicPeriodResponse, 0, len(tour.AvailablePeriods))
		for _, period := range tour.AvailablePeriods {
			// Only include published periods
			if period.Status == "publish" {
				publicPeriod := &dto.PublicPeriodResponse{
					ID:                 period.ID,
					TourID:             period.TourID,
					StartDate:          period.StartDate,
					EndDate:            period.EndDate,
					Seat:               period.Seat,
					Booking:            period.Booking,
					HeldSeats:          period.HeldSeats,
					AvailableSeats:     period.AvailableSeats,
					Deposit:            period.Deposit,
					DepositType:        period.DepositType,
					Airline:            period.Airline,
					Airport:            period.Airport,
					Note:               period.Note,
					Sector1:            period.Sector1,
					Sector2:            period.Sector2,
					RateType:           period.RateType,
					RateItems:          period.RateItems,
					RoomTypes:          period.RoomTypes,
					AdditionalServices: period.AdditionalServices,
					Files:              period.Files,
					Status:             period.Status,
					CreatedAt:          period.CreatedAt,
					UpdatedAt:          period.UpdatedAt,
				}
				publicPeriods = append(publicPeriods, publicPeriod)
			}
		}
		publicTour.AvailablePeriods = publicPeriods
	}

	return publicTour
}
