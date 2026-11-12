package services

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
)

type BannerService interface {
	// CRUD Operations
	CreateBanner(ctx context.Context, req dto.CreateBannerRequest) (*dto.BannerResponse, error)
	UpdateBanner(ctx context.Context, bannerID uuid.UUID, req dto.UpdateBannerRequest) (*dto.BannerResponse, error)
	DeleteBanner(ctx context.Context, bannerID uuid.UUID) error
	GetBanner(ctx context.Context, bannerID uuid.UUID) (*dto.BannerResponse, error)
	GetBannerByID(ctx context.Context, bannerID uuid.UUID) (*dto.BannerResponse, error) // Alias for GetBanner
	ListBanners(ctx context.Context, req dto.BannerListRequest) (*dto.BannerListResponse, error)

	// Status Management
	UpdateStatus(ctx context.Context, req dto.UpdateBannerStatusRequest) (*dto.BannerResponse, error)
	UpdateBannerStatus(ctx context.Context, bannerID uuid.UUID, req dto.UpdateBannerStatusRequest) (*dto.BannerResponse, error)
	PublishBanner(ctx context.Context, bannerID uuid.UUID) (*dto.BannerResponse, error)
	UnpublishBanner(ctx context.Context, bannerID uuid.UUID) (*dto.BannerResponse, error)

	// Order Management
	UpdateOrder(ctx context.Context, req dto.UpdateBannerOrderRequest) error
	UpdateBannerOrder(ctx context.Context, bannerID uuid.UUID, req dto.UpdateBannerOrderRequest) (*dto.BannerResponse, error)

	// Public Operations
	GetPublishedBanners(ctx context.Context, offset, limit int) ([]*dto.BannerResponse, error)
	GetActiveBanners(ctx context.Context, position string) ([]*dto.BannerResponse, error)
	GetRecentPublishedBanners(ctx context.Context, limit int) ([]*dto.BannerResponse, error)

	// Statistics
	GetBannerStats(ctx context.Context) (*dto.BannerStatsResponse, error)
}