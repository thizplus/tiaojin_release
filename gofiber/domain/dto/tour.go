package dto

import "github.com/google/uuid"

type TourResponse struct {
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

type CreateTourRequest struct {
	Title     string     `json:"title" validate:"required,max=300"`
	Code      string     `json:"code" validate:"required,max=20,uppercase"`
	CityID    *uuid.UUID `json:"city_id,omitempty"`
	Detail    *string    `json:"detail,omitempty"`
	Thumbnail *string    `json:"thumbnail,omitempty"`
	FilePDF   *string    `json:"file_pdf,omitempty"`
	FileWord  *string    `json:"file_word,omitempty"`
	Status    *string    `json:"status,omitempty" validate:"omitempty,oneof=publish draft"`
}

type UpdateTourRequest struct {
	Title     *string    `json:"title,omitempty" validate:"omitempty,max=300"`
	CityID    *uuid.UUID `json:"city_id,omitempty"`
	Detail    *string    `json:"detail,omitempty"`
	Thumbnail *string    `json:"thumbnail,omitempty"`
	FilePDF   *string    `json:"file_pdf,omitempty"`
	FileWord  *string    `json:"file_word,omitempty"`
	Status    *string    `json:"status,omitempty" validate:"omitempty,oneof=publish draft"`
}

type UploadTourFileRequest struct {
	TourID   uuid.UUID `json:"tour_id" validate:"required"`
	FileData []byte    `json:"file_data" validate:"required"`
	FileName string    `json:"file_name" validate:"required"`
	FileType string    `json:"file_type" validate:"required,oneof=application/pdf application/msword application/vnd.openxmlformats-officedocument.wordprocessingml.document"`
}

type TourListRequest struct {
	CityID      *uuid.UUID `json:"city_id,omitempty"`
	Status      *string    `json:"status,omitempty" validate:"omitempty,oneof=publish draft"`
	Search      *string    `json:"search,omitempty"`
	WithPeriods bool       `json:"with_periods,omitempty"`
	Offset      int        `json:"offset" validate:"min=0"`
	Limit       int        `json:"limit" validate:"min=1,max=100"`
}

type TourListResponse struct {
	Tours      []*TourResponse `json:"tours"`
	Pagination PaginationMeta  `json:"pagination"`
}

type TourWithPeriodsResponse struct {
	*TourResponse
	AvailablePeriods []*PeriodResponse `json:"available_periods,omitempty"`
	MinStartDate     *string           `json:"min_start_date,omitempty"`
}

type TourWithPeriodsListResponse struct {
	Tours      []*TourWithPeriodsResponse `json:"tours"`
	Pagination PaginationMeta             `json:"pagination"`
}

type PeriodSummaryResponse struct {
	ID             uuid.UUID                     `json:"id"`
	StartDate      string                        `json:"start_date"`
	EndDate        string                        `json:"end_date"`
	Seat           int                           `json:"seat"`
	Booking        int                           `json:"booking"`
	HeldSeats      int                           `json:"held_seats"`
	AvailableSeats int                           `json:"available_seats"`
	Deposit        float64                       `json:"deposit"`
	DepositType    string                        `json:"deposit_type"`
	Files          []*FileAttachmentResponse     `json:"files"`
}

type TourStatsResponse struct {
	TotalTours       int64 `json:"total_tours"`
	PublishedTours   int64 `json:"published_tours"`
	DraftTours       int64 `json:"draft_tours"`
	ToursWithPeriods int64 `json:"tours_with_periods"`
}

type UpdateTourStatusRequest struct {
	ID     uuid.UUID `json:"id" validate:"required"`
	Status string    `json:"status" validate:"required,oneof=publish draft"`
}

