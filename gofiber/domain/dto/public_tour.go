package dto

import "github.com/google/uuid"

// PublicTourResponse - Tour response for public API (no commission info)
type PublicTourResponse struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Code      string     `json:"code"`
	CityID    *uuid.UUID `json:"city_id"`
	CityTitle *string    `json:"city_title,omitempty"`
	Detail    *string    `json:"detail"`
	Thumbnail *string    `json:"thumbnail"`
	FilePDF   *string    `json:"file_pdf"`
	FileWord  *string    `json:"file_word"`
	Status    string     `json:"status"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

// PublicPeriodResponse - Period response for public API (no commission info)
type PublicPeriodResponse struct {
	ID                 uuid.UUID                    `json:"id"`
	TourID             uuid.UUID                    `json:"tour_id"`
	StartDate          string                       `json:"start_date"`
	EndDate            string                       `json:"end_date"`
	Seat               int                          `json:"seat"`
	Booking            int                          `json:"booking"`
	HeldSeats          int                          `json:"held_seats"`
	AvailableSeats     int                          `json:"available_seats"`
	Deposit            float64                      `json:"deposit"`
	DepositType        string                       `json:"deposit_type"`
	Airline            string                       `json:"airline"`
	Airport            string                       `json:"airport"`
	Note               *string                      `json:"note"`
	Sector1            []*FlightSectorResponse      `json:"sector_1"`
	Sector2            []*FlightSectorResponse      `json:"sector_2"`
	RateType           string                       `json:"rate_type"`
	RateItems          []*RateItemResponse          `json:"rate_items"`
	RoomTypes          []*RoomTypeResponse          `json:"room_types"`
	AdditionalServices []*AdditionalServiceResponse `json:"additional_services"`
	Files              []*FileAttachmentResponse    `json:"files"`
	Status             string                       `json:"status"`
	CreatedAt          string                       `json:"created_at"`
	UpdatedAt          string                       `json:"updated_at"`
}

// PublicTourWithPeriodsResponse - Tour with periods for public API
type PublicTourWithPeriodsResponse struct {
	*PublicTourResponse
	AvailablePeriods []*PublicPeriodResponse `json:"available_periods,omitempty"`
	MinStartDate     *string                 `json:"min_start_date,omitempty"`
}

// PublicTourListResponse - List response for public API
type PublicTourListResponse struct {
	Tours      []*PublicTourWithPeriodsResponse `json:"tours"`
	Pagination PaginationMeta                   `json:"pagination"`
}

// PublicTourListRequest - List request for public API
type PublicTourListRequest struct {
	CityID *uuid.UUID `json:"city_id,omitempty"`
	Search *string    `json:"search,omitempty"`
	Offset int        `json:"offset" validate:"min=0"`
	Limit  int        `json:"limit" validate:"min=1,max=100"`
}
