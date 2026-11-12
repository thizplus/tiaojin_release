package dto

import (
	"time"

	"github.com/google/uuid"
)

type PeriodResponse struct {
	ID             uuid.UUID `json:"id"`
	TourID         uuid.UUID `json:"tour_id"`
	StartDate      string    `json:"start_date"`
	EndDate        string    `json:"end_date"`
	Seat           int       `json:"seat"`
	Booking        int       `json:"booking"`
	HeldSeats      int       `json:"held_seats"`
	AvailableSeats int       `json:"available_seats"`
	Deposit        float64   `json:"deposit"`
	DepositType    string    `json:"deposit_type"`
	//SingleStay      float64                `json:"single_stay"`
	Airline            string                       `json:"airline"`
	Airport            string                       `json:"airport"`
	Note               *string                      `json:"note"`
	Sector1            []*FlightSectorResponse      `json:"sector_1"`
	Sector2            []*FlightSectorResponse      `json:"sector_2"`
	RateType           string                       `json:"rate_type"`
	RateItems          []*RateItemResponse          `json:"rate_items"`
	CommissionType     string                       `json:"commission_type"`
	CommissionItems    []*CommissionItemResponse    `json:"commission_items"`
	RoomTypes          []*RoomTypeResponse          `json:"room_types"`
	AdditionalServices []*AdditionalServiceResponse `json:"additional_services"`
	Files              []*FileAttachmentResponse    `json:"files"`
	Status             string                       `json:"status"`
	ServiceFeePerPax   float64                      `json:"service_fee_per_pax"`
	CreatedAt          string                       `json:"created_at"`
	UpdatedAt          string                       `json:"updated_at"`

	// Expiration Override
	ExpirationOverride bool    `json:"expiration_override"`
	ExpirationMinutes  *int    `json:"expiration_minutes,omitempty"`
	ExpirationHours    *int    `json:"expiration_hours,omitempty"`
	ExpirationDays     *int    `json:"expiration_days,omitempty"`
	ExpirationMessage  *string `json:"expiration_message,omitempty"`
	ExpirationDisplay  string  `json:"expiration_display,omitempty"` // Human-readable duration

	// Waiting List
	HasWaitingList bool `json:"has_waiting_list"`

	// Relations
	TourDetails *TourResponse `json:"tour_details,omitempty"`
}

type CreatePeriodRequest struct {
	TourID      uuid.UUID `json:"tour_id" validate:"required"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`
	Seat        int       `json:"seat" validate:"required,min=1,max=999"`
	Deposit     float64   `json:"deposit" validate:"required,min=0"`
	DepositType string    `json:"deposit_type" validate:"required,oneof=percent per_unit"`
	//	SingleStay         float64                           `json:"single_stay" validate:"min=0"`
	Airline            string                            `json:"airline" validate:"required,max=100"`
	Airport            string                            `json:"airport" validate:"required,max=10"`
	Note               *string                           `json:"note,omitempty"`
	Sector1            []*CreateFlightSectorRequest      `json:"sector_1"`
	Sector2            []*CreateFlightSectorRequest      `json:"sector_2"`
	RateType           string                            `json:"rate_type" validate:"required,oneof=regular sale"`
	RateItems          []*CreateRateItemRequest          `json:"rate_items" validate:"required,min=1"`
	CommissionType     string                            `json:"commission_type" validate:"required,oneof=regular sale"`
	CommissionItems    []*CreateCommissionItemRequest    `json:"commission_items" validate:"required,min=1"`
	RoomTypes          []*CreateRoomTypeRequest          `json:"room_types,omitempty"`
	AdditionalServices []*CreateAdditionalServiceRequest `json:"additional_services,omitempty"`
	Status             *string                           `json:"status,omitempty" validate:"omitempty,oneof=publish draft"`
	ServiceFeePerPax   float64                           `json:"service_fee_per_pax" validate:"min=0"`

	// Expiration Override - สำหรับทัวร์ไฟไหม้
	ExpirationOverride *bool   `json:"expiration_override,omitempty"`
	ExpirationMinutes  *int    `json:"expiration_minutes,omitempty" validate:"omitempty,min=1,max=10080"`
	ExpirationHours    *int    `json:"expiration_hours,omitempty" validate:"omitempty,min=1,max=168"`
	ExpirationDays     *int    `json:"expiration_days,omitempty" validate:"omitempty,min=1,max=30"`
	ExpirationMessage  *string `json:"expiration_message,omitempty" validate:"omitempty,max=255"`
}

type UpdatePeriodRequest struct {
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Seat        *int       `json:"seat,omitempty" validate:"omitempty,min=1,max=999"`
	Deposit     *float64   `json:"deposit,omitempty" validate:"omitempty,min=0"`
	DepositType *string    `json:"deposit_type,omitempty" validate:"omitempty,oneof=percent per_unit"`
	//SingleStay         *float64                          `json:"single_stay,omitempty" validate:"omitempty,min=0"`
	Airline            *string                           `json:"airline,omitempty" validate:"omitempty,max=100"`
	Airport            *string                           `json:"airport,omitempty" validate:"omitempty,max=10"`
	Note               *string                           `json:"note,omitempty"`
	Sector1            []*CreateFlightSectorRequest      `json:"sector_1,omitempty"`
	Sector2            []*CreateFlightSectorRequest      `json:"sector_2,omitempty"`
	RateType           *string                           `json:"rate_type,omitempty" validate:"omitempty,oneof=regular sale"`
	RateItems          []*CreateRateItemRequest          `json:"rate_items,omitempty"`
	CommissionType     *string                           `json:"commission_type,omitempty" validate:"omitempty,oneof=regular sale"`
	CommissionItems    []*CreateCommissionItemRequest    `json:"commission_items,omitempty"`
	RoomTypes          []*CreateRoomTypeRequest          `json:"room_types,omitempty"`
	AdditionalServices []*CreateAdditionalServiceRequest `json:"additional_services,omitempty"`
	Files              []*FileAttachmentRequest          `json:"files,omitempty"`
	Status             *string                           `json:"status,omitempty" validate:"omitempty,oneof=publish draft"`
	ServiceFeePerPax   *float64                          `json:"service_fee_per_pax,omitempty" validate:"omitempty,min=0"`

	// Expiration Override - สำหรับทัวร์ไฟไหม้
	ExpirationOverride *bool   `json:"expiration_override,omitempty"`
	ExpirationMinutes  *int    `json:"expiration_minutes,omitempty" validate:"omitempty,min=1,max=10080"`
	ExpirationHours    *int    `json:"expiration_hours,omitempty" validate:"omitempty,min=1,max=168"`
	ExpirationDays     *int    `json:"expiration_days,omitempty" validate:"omitempty,min=1,max=30"`
	ExpirationMessage  *string `json:"expiration_message,omitempty" validate:"omitempty,max=255"`
}

type CreateFlightSectorRequest struct {
	Date   string `json:"date" validate:"required"`
	Flight string `json:"flight" validate:"required"`
	Route  string `json:"route" validate:"required"`
}

type CreateRateItemRequest struct {
	ID        string  `json:"id,omitempty"`
	Title     string  `json:"title" validate:"required,max=100"`
	Regular   float64 `json:"regular" validate:"required,min=0"`
	Sale      float64 `json:"sale" validate:"required,min=0"`
	Exclude   bool    `json:"exclude"`
	SortOrder *int    `json:"sort_order,omitempty"`
}

type CreateCommissionItemRequest struct {
	Regular float64 `json:"regular" validate:"required,min=0"`
	Sale    float64 `json:"sale" validate:"required,min=0"`
}

type CreateRoomTypeRequest struct {
	ID        string  `json:"id,omitempty"`
	Title     string  `json:"title" validate:"required,max=100"`
	Price     float64 `json:"price" validate:"required,min=0"`
	SortOrder *int    `json:"sort_order,omitempty"`
}

type CreateAdditionalServiceRequest struct {
	ID        string  `json:"id,omitempty"`
	Title     string  `json:"title" validate:"required,max=100"`
	Price     float64 `json:"price" validate:"required,min=0"`
	RateType  string  `json:"rate_type" validate:"required,oneof=per_unit per_booking"`
	UnitName  string  `json:"unit_name" validate:"required"`
	SortOrder *int    `json:"sort_order,omitempty"`
}

type FlightSectorResponse struct {
	Date   string `json:"date"`
	Flight string `json:"flight"`
	Route  string `json:"route"`
}

type RateItemResponse struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Regular   float64 `json:"regular"`
	Sale      float64 `json:"sale"`
	Exclude   bool    `json:"exclude"`
	SortOrder int     `json:"sort_order"`
}

type CommissionItemResponse struct {
	Regular float64 `json:"regular"`
	Sale    float64 `json:"sale"`
}

type RoomTypeResponse struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	SortOrder int     `json:"sort_order"`
}

type AdditionalServiceResponse struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	RateType  string  `json:"rate_type"`
	UnitName  string  `json:"unit_name"`
	SortOrder int     `json:"sort_order"`
}

type FileAttachmentRequest struct {
	Name string `json:"name" validate:"required,max=255"`
	Type string `json:"type" validate:"required,max=100"`
	URL  string `json:"url" validate:"required,url,max=500"`
}

type FileAttachmentResponse struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type PeriodListRequest struct {
	TourID    *uuid.UUID `json:"tourId,omitempty"`
	Status    *string    `json:"status,omitempty" validate:"omitempty,oneof=publish draft"`
	DateFrom  *time.Time `json:"dateFrom,omitempty"`
	DateTo    *time.Time `json:"dateTo,omitempty"`
	Available *bool      `json:"available,omitempty"`
	Upcoming  *bool      `json:"upcoming,omitempty"`
	Offset    int        `json:"offset" validate:"min=0"`
	Limit     int        `json:"limit" validate:"min=1,max=100"`
}

type PeriodListResponse struct {
	Periods    []*PeriodResponse `json:"periods"`
	Pagination PaginationMeta    `json:"pagination"`
}

type AvailabilityCheckResponse struct {
	Available      bool `json:"available"`
	RequestedSeats int  `json:"requested_seats"`
	AvailableSeats int  `json:"available_seats"`
	TotalSeats     int  `json:"total_seats"`
	BookedSeats    int  `json:"booked_seats"`
}

type PeriodStatsResponse struct {
	TotalPeriods       int64   `json:"total_periods"`
	PublishedPeriods   int64   `json:"published_periods"`
	AvailablePeriods   int64   `json:"available_periods"`
	FullyBookedPeriods int64   `json:"fully_booked_periods"`
	TotalSeats         int64   `json:"total_seats"`
	BookedSeats        int64   `json:"booked_seats"`
	AverageOccupancy   float64 `json:"average_occupancy"`
}

type UpdatePeriodStatusRequest struct {
	ID     uuid.UUID `json:"id" validate:"required"`
	Status string    `json:"status" validate:"required,oneof=publish draft"`
}

type PeriodAvailabilityResponse struct {
	PeriodID       uuid.UUID `json:"periodId"`
	TotalSeats     int       `json:"totalSeats"`
	BookedSeats    int       `json:"bookedSeats"`
	AvailableSeats int       `json:"availableSeats"`
	IsAvailable    bool      `json:"isAvailable"`
}

type PeriodBookingStatsResponse struct {
	PeriodID         uuid.UUID      `json:"period_id"`
	TotalSeats       int            `json:"total_seats"`
	TotalBookings    int            `json:"total_bookings"`
	BookedPassengers int            `json:"booked_passengers"`
	HeldSeats        int            `json:"held_seats"`
	AvailableSeats   int            `json:"available_seats"`
	BookingsByStatus map[string]int `json:"bookings_by_status"`
	IsAvailable      bool           `json:"is_available"`
}
