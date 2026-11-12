package dto

import "github.com/google/uuid"

// Request DTOs
type CreateCityRequest struct {
	Title  string  `json:"title" validate:"required,min=1,max=200"`
	Detail *string `json:"detail,omitempty" validate:"omitempty,max=5000"`
	Image  string  `json:"image" validate:"required,max=500"`
	Status string  `json:"status" validate:"required,oneof=publish draft"`
}

type UpdateCityRequest struct {
	Title  *string `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Detail *string `json:"detail,omitempty" validate:"omitempty,max=5000"`
	Image  *string `json:"image,omitempty" validate:"omitempty,max=500"`
	Status *string `json:"status,omitempty" validate:"omitempty,oneof=publish draft"`
}

type UpdateCityStatusRequest struct {
	CityID uuid.UUID `json:"city_id" validate:"required"`
	Status string    `json:"status" validate:"required,oneof=publish draft"`
}

type CityListRequest struct {
	Status *string `json:"status,omitempty" validate:"omitempty,oneof=publish draft"`
	Search *string `json:"search,omitempty" validate:"omitempty,min=2,max=50"`
	Stat   *bool   `json:"stat,omitempty"`
	Offset int     `json:"offset" validate:"min=0"`
	Limit  int     `json:"limit" validate:"min=1,max=100"`
}

// Response DTOs
type CityResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Detail    *string   `json:"detail"`
	Image     string    `json:"image"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`

	// Statistics (only populated when stat=true)
	ActiveTours   *int `json:"active_tours,omitempty"`
	ActivePeriods *int `json:"active_periods,omitempty"`
}

type CityListResponse struct {
	Cities     []*CityResponse `json:"cities"`
	Pagination PaginationMeta  `json:"pagination"`
}