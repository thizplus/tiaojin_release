package dto

import "github.com/google/uuid"

// Request DTOs
type CreateBannerRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=200"`
	Image        string  `json:"image" validate:"required,max=500"`
	LinkURL      *string `json:"link_url,omitempty" validate:"omitempty,max=500"`
	Alt          *string `json:"alt,omitempty" validate:"omitempty,max=200"`
	DisplayOrder *int    `json:"display_order,omitempty" validate:"omitempty,min=0"`
	Position     string  `json:"position" validate:"omitempty,oneof=public internal"`
	Status       string  `json:"status" validate:"required,oneof=publish draft"`
}

type UpdateBannerRequest struct {
	Title        *string `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Image        *string `json:"image,omitempty" validate:"omitempty,max=500"`
	LinkURL      *string `json:"link_url,omitempty" validate:"omitempty,max=500"`
	Alt          *string `json:"alt,omitempty" validate:"omitempty,max=200"`
	DisplayOrder *int    `json:"display_order,omitempty" validate:"omitempty,min=0"`
	Position     *string `json:"position,omitempty" validate:"omitempty,oneof=public internal"`
	Status       *string `json:"status,omitempty" validate:"omitempty,oneof=publish draft"`
}

type UpdateBannerStatusRequest struct {
	BannerID uuid.UUID `json:"banner_id" validate:"required"`
	Status   string    `json:"status" validate:"required,oneof=publish draft"`
}

type UpdateBannerOrderRequest struct {
	DisplayOrder int               `json:"display_order" validate:"required,min=0"`
	BannerOrders []BannerOrderItem `json:"banner_orders,omitempty" validate:"omitempty,min=1"`
}

type UpdateBannerBulkOrderRequest struct {
	BannerOrders []BannerOrderItem `json:"banner_orders" validate:"required,min=1"`
}

type BannerOrderItem struct {
	BannerID uuid.UUID `json:"banner_id" validate:"required"`
	Order    int       `json:"order" validate:"required,min=0"`
}

type BannerListRequest struct {
	Search   *string `json:"search,omitempty" validate:"omitempty,min=1,max=100"`
	Status   *string `json:"status,omitempty" validate:"omitempty,oneof=publish draft"`
	Position *string `json:"position,omitempty"`
	Offset   int     `json:"offset" validate:"min=0"`
	Limit    int     `json:"limit" validate:"min=1,max=100"`
}

// Response DTOs
type BannerResponse struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Image        string    `json:"image"`
	LinkURL      *string   `json:"link_url"`
	Alt          *string   `json:"alt"`
	DisplayOrder int       `json:"display_order"`
	Position     string    `json:"position"`
	Status       string    `json:"status"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

type BannerListResponse struct {
	Banners    []*BannerResponse `json:"banners"`
	Pagination PaginationMeta    `json:"pagination"`
}

type BannerStatsResponse struct {
	TotalBanners     int64 `json:"total_banners"`
	PublishedBanners int64 `json:"published_banners"`
	DraftBanners     int64 `json:"draft_banners"`
}