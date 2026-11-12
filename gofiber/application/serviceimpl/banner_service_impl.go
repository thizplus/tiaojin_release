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

type BannerServiceImpl struct {
	bannerRepo repositories.BannerRepository
}

func NewBannerService(bannerRepo repositories.BannerRepository) services.BannerService {
	return &BannerServiceImpl{
		bannerRepo: bannerRepo,
	}
}

// CRUD Operations
func (s *BannerServiceImpl) CreateBanner(ctx context.Context, req dto.CreateBannerRequest) (*dto.BannerResponse, error) {
	// Determine position (default: internal)
	position := models.BannerPositionInternal
	if req.Position == "public" {
		position = models.BannerPositionPublic
	}

	// Create domain model
	banner := &models.Banner{
		Title:        req.Title,
		Image:        req.Image,
		LinkURL:      req.LinkURL,
		Alt:          req.Alt,
		DisplayOrder: getDisplayOrderOrDefault(req.DisplayOrder, 0),
		Position:     position,
		Status:       models.BannerStatus(req.Status), // Use status from request
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Create banner
	if err := s.bannerRepo.Create(ctx, banner); err != nil {
		log.Printf("Error creating banner: %v", err)
		return nil, fmt.Errorf("failed to create banner: %w", err)
	}

	return s.toBannerResponse(banner), nil
}

func (s *BannerServiceImpl) UpdateBanner(ctx context.Context, bannerID uuid.UUID, req dto.UpdateBannerRequest) (*dto.BannerResponse, error) {
	// Get existing banner
	banner, err := s.bannerRepo.GetByID(ctx, bannerID)
	if err != nil {
		log.Printf("Error retrieving banner: %v", err)
		return nil, fmt.Errorf("failed to get banner: %w", err)
	}

	if banner == nil {
		return nil, apperrors.ErrBannerNotFound
	}

	// Update fields if provided
	if req.Title != nil {
		banner.Title = *req.Title
	}
	if req.Image != nil {
		banner.Image = *req.Image
	}
	if req.LinkURL != nil {
		banner.LinkURL = req.LinkURL
	}
	if req.Alt != nil {
		banner.Alt = req.Alt
	}
	if req.DisplayOrder != nil {
		banner.DisplayOrder = *req.DisplayOrder
	}
	if req.Position != nil {
		banner.Position = models.BannerPosition(*req.Position)
	}
	if req.Status != nil {
		banner.Status = models.BannerStatus(*req.Status)
	}

	banner.UpdatedAt = time.Now()

	// Update banner
	if err := s.bannerRepo.Update(ctx, bannerID, banner); err != nil {
		log.Printf("Error updating banner: %v", err)
		return nil, fmt.Errorf("failed to update banner: %w", err)
	}

	return s.toBannerResponse(banner), nil
}

func (s *BannerServiceImpl) DeleteBanner(ctx context.Context, bannerID uuid.UUID) error {
	// Check if banner exists
	banner, err := s.bannerRepo.GetByID(ctx, bannerID)
	if err != nil {
		log.Printf("Error retrieving banner: %v", err)
		return fmt.Errorf("failed to get banner: %w", err)
	}

	if banner == nil {
		return apperrors.ErrBannerNotFound
	}

	// Delete banner
	if err := s.bannerRepo.Delete(ctx, bannerID); err != nil {
		log.Printf("Error deleting banner: %v", err)
		return fmt.Errorf("failed to delete banner: %w", err)
	}

	return nil
}

func (s *BannerServiceImpl) GetBanner(ctx context.Context, bannerID uuid.UUID) (*dto.BannerResponse, error) {
	banner, err := s.bannerRepo.GetByID(ctx, bannerID)
	if err != nil {
		log.Printf("Error retrieving banner: %v", err)
		return nil, fmt.Errorf("failed to get banner: %w", err)
	}

	if banner == nil {
		return nil, apperrors.ErrBannerNotFound
	}

	return s.toBannerResponse(banner), nil
}

func (s *BannerServiceImpl) GetBannerByID(ctx context.Context, bannerID uuid.UUID) (*dto.BannerResponse, error) {
	// This is an alias for GetBanner to maintain compatibility with handlers
	return s.GetBanner(ctx, bannerID)
}

func (s *BannerServiceImpl) ListBanners(ctx context.Context, req dto.BannerListRequest) (*dto.BannerListResponse, error) {
	var banners []*models.Banner
	var total int64
	var err error

	// Helper: filter by position if specified
	filterByPosition := func(items []*models.Banner) []*models.Banner {
		if req.Position == nil || *req.Position == "" {
			return items
		}
		filtered := make([]*models.Banner, 0)
		for _, b := range items {
			if string(b.Position) == *req.Position {
				filtered = append(filtered, b)
			}
		}
		return filtered
	}

	// Handle search if provided
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		banners, err = s.bannerRepo.SearchByTitle(ctx, strings.TrimSpace(*req.Search), req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error searching banners: %v", err)
			return nil, fmt.Errorf("failed to search banners: %w", err)
		}
		banners = filterByPosition(banners)
		total = int64(len(banners))
	} else if req.Status != nil && *req.Status == "publish" {
		// Get published banners
		banners, err = s.bannerRepo.GetPublished(ctx, req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error getting published banners: %v", err)
			return nil, fmt.Errorf("failed to get published banners: %w", err)
		}

		// Get published count
		total, err = s.bannerRepo.GetPublishedCount(ctx)
		if err != nil {
			log.Printf("Error counting published banners: %v", err)
			return nil, fmt.Errorf("failed to count published banners: %w", err)
		}
	} else {
		// List all banners (with optional draft filter)
		banners, err = s.bannerRepo.List(ctx, req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error listing banners: %v", err)
			return nil, fmt.Errorf("failed to list banners: %w", err)
		}

		// Get total count
		total, err = s.bannerRepo.Count(ctx)
		if err != nil {
			log.Printf("Error counting banners: %v", err)
			return nil, fmt.Errorf("failed to count banners: %w", err)
		}

		// Filter by draft status if specified
		if req.Status != nil && *req.Status == "draft" {
			filteredBanners := make([]*models.Banner, 0)
			for _, banner := range banners {
				if banner.Status == models.BannerStatusDraft {
					filteredBanners = append(filteredBanners, banner)
				}
			}
			banners = filteredBanners
		}
	}

	// Filter by position
	banners = filterByPosition(banners)

	// Convert to response DTOs
	responseBanners := make([]*dto.BannerResponse, len(banners))
	for i, banner := range banners {
		responseBanners[i] = s.toBannerResponse(banner)
	}

	return &dto.BannerListResponse{
		Banners:    responseBanners,
		Pagination: utils.CalculatePagination(total, req.Offset, req.Limit),
	}, nil
}

// Status Management
func (s *BannerServiceImpl) UpdateStatus(ctx context.Context, req dto.UpdateBannerStatusRequest) (*dto.BannerResponse, error) {
	// Business validation - validate status
	if req.Status != "publish" && req.Status != "draft" {
		return nil, apperrors.ErrInvalidBannerStatus
	}

	// Get existing banner
	banner, err := s.bannerRepo.GetByID(ctx, req.BannerID)
	if err != nil {
		log.Printf("Error retrieving banner: %v", err)
		return nil, fmt.Errorf("failed to get banner: %w", err)
	}

	if banner == nil {
		return nil, apperrors.ErrBannerNotFound
	}

	// Update status
	status := models.BannerStatus(req.Status)
	if err := s.bannerRepo.UpdateStatus(ctx, req.BannerID, status); err != nil {
		log.Printf("Error updating banner status: %v", err)
		return nil, fmt.Errorf("failed to update banner status: %w", err)
	}

	banner.Status = status
	banner.UpdatedAt = time.Now()

	return s.toBannerResponse(banner), nil
}

func (s *BannerServiceImpl) PublishBanner(ctx context.Context, bannerID uuid.UUID) (*dto.BannerResponse, error) {
	// Get existing banner
	banner, err := s.bannerRepo.GetByID(ctx, bannerID)
	if err != nil {
		log.Printf("Error retrieving banner: %v", err)
		return nil, fmt.Errorf("failed to get banner: %w", err)
	}

	if banner == nil {
		return nil, apperrors.ErrBannerNotFound
	}

	// Publish banner
	if err := s.bannerRepo.Publish(ctx, bannerID); err != nil {
		log.Printf("Error publishing banner: %v", err)
		return nil, fmt.Errorf("failed to publish banner: %w", err)
	}

	banner.Status = models.BannerStatusPublish
	banner.UpdatedAt = time.Now()

	return s.toBannerResponse(banner), nil
}

func (s *BannerServiceImpl) UnpublishBanner(ctx context.Context, bannerID uuid.UUID) (*dto.BannerResponse, error) {
	// Get existing banner
	banner, err := s.bannerRepo.GetByID(ctx, bannerID)
	if err != nil {
		log.Printf("Error retrieving banner: %v", err)
		return nil, fmt.Errorf("failed to get banner: %w", err)
	}

	if banner == nil {
		return nil, apperrors.ErrBannerNotFound
	}

	// Unpublish banner
	if err := s.bannerRepo.Unpublish(ctx, bannerID); err != nil {
		log.Printf("Error unpublishing banner: %v", err)
		return nil, fmt.Errorf("failed to unpublish banner: %w", err)
	}

	banner.Status = models.BannerStatusDraft
	banner.UpdatedAt = time.Now()

	return s.toBannerResponse(banner), nil
}

// Order Management
func (s *BannerServiceImpl) UpdateOrder(ctx context.Context, req dto.UpdateBannerOrderRequest) error {
	// Business validation - check if list is not empty
	if len(req.BannerOrders) == 0 {
		return apperrors.ErrEmptyBannerList
	}

	// Business validation - check for duplicate IDs
	seen := make(map[uuid.UUID]bool)
	for _, order := range req.BannerOrders {
		if seen[order.BannerID] {
			return apperrors.ErrDuplicateBannerID
		}
		seen[order.BannerID] = true
	}

	// Business validation - check if all banners exist
	for _, order := range req.BannerOrders {
		banner, err := s.bannerRepo.GetByID(ctx, order.BannerID)
		if err != nil {
			log.Printf("Error checking banner existence: %v", err)
			return fmt.Errorf("failed to validate banner: %w", err)
		}
		if banner == nil {
			return fmt.Errorf("banner with ID %s not found", order.BannerID)
		}
	}

	// TODO: Implement order management
	// This would typically involve updating an order field in the database
	// For now, just log the operation
	log.Printf("Updating banner order for %d banners", len(req.BannerOrders))

	return nil
}

// Public Operations
func (s *BannerServiceImpl) GetPublishedBanners(ctx context.Context, offset, limit int) ([]*dto.BannerResponse, error) {
	banners, err := s.bannerRepo.GetPublished(ctx, offset, limit)
	if err != nil {
		log.Printf("Error getting published banners: %v", err)
		return nil, fmt.Errorf("failed to get published banners: %w", err)
	}

	// Convert to response DTOs
	responseBanners := make([]*dto.BannerResponse, len(banners))
	for i, banner := range banners {
		responseBanners[i] = s.toBannerResponse(banner)
	}

	return responseBanners, nil
}

func (s *BannerServiceImpl) GetRecentPublishedBanners(ctx context.Context, limit int) ([]*dto.BannerResponse, error) {
	banners, err := s.bannerRepo.GetRecentPublished(ctx, limit)
	if err != nil {
		log.Printf("Error getting recent published banners: %v", err)
		return nil, fmt.Errorf("failed to get recent published banners: %w", err)
	}

	// Convert to response DTOs
	responseBanners := make([]*dto.BannerResponse, len(banners))
	for i, banner := range banners {
		responseBanners[i] = s.toBannerResponse(banner)
	}

	return responseBanners, nil
}

// Statistics
func (s *BannerServiceImpl) GetBannerStats(ctx context.Context) (*dto.BannerStatsResponse, error) {
	// Get total banners count
	totalBanners, err := s.bannerRepo.Count(ctx)
	if err != nil {
		log.Printf("Error getting banner stats: %v", err)
		return nil, fmt.Errorf("failed to get banner stats: %w", err)
	}

	// Get published banners count
	publishedBanners, err := s.bannerRepo.GetPublishedCount(ctx)
	if err != nil {
		log.Printf("Error getting published banner stats: %v", err)
		return nil, fmt.Errorf("failed to get banner stats: %w", err)
	}

	draftBanners := totalBanners - publishedBanners

	return &dto.BannerStatsResponse{
		TotalBanners:     totalBanners,
		PublishedBanners: publishedBanners,
		DraftBanners:     draftBanners,
	}, nil
}

// Additional methods to match handler expectations
func (s *BannerServiceImpl) UpdateBannerStatus(ctx context.Context, bannerID uuid.UUID, req dto.UpdateBannerStatusRequest) (*dto.BannerResponse, error) {
	// Override the bannerID in the request with the one from URL
	req.BannerID = bannerID
	return s.UpdateStatus(ctx, req)
}

func (s *BannerServiceImpl) UpdateBannerOrder(ctx context.Context, bannerID uuid.UUID, req dto.UpdateBannerOrderRequest) (*dto.BannerResponse, error) {
	// Get existing banner
	banner, err := s.bannerRepo.GetByID(ctx, bannerID)
	if err != nil {
		log.Printf("Error retrieving banner: %v", err)
		return nil, fmt.Errorf("failed to get banner: %w", err)
	}

	if banner == nil {
		return nil, apperrors.ErrBannerNotFound
	}

	// Update display order
	banner.DisplayOrder = req.DisplayOrder
	banner.UpdatedAt = time.Now()

	// Update banner in repository
	if err := s.bannerRepo.Update(ctx, bannerID, banner); err != nil {
		log.Printf("Error updating banner order: %v", err)
		return nil, fmt.Errorf("failed to update banner order: %w", err)
	}

	return s.toBannerResponse(banner), nil
}

func (s *BannerServiceImpl) GetActiveBanners(ctx context.Context, position string) ([]*dto.BannerResponse, error) {
	// Get published banners
	banners, err := s.bannerRepo.GetPublished(ctx, 0, 100)
	if err != nil {
		log.Printf("Error getting active banners: %v", err)
		return nil, fmt.Errorf("failed to get active banners: %w", err)
	}

	// Filter by position if specified
	if position != "" {
		filtered := make([]*models.Banner, 0)
		for _, banner := range banners {
			if string(banner.Position) == position {
				filtered = append(filtered, banner)
			}
		}
		banners = filtered
	}

	// Convert to response DTOs
	responseBanners := make([]*dto.BannerResponse, len(banners))
	for i, banner := range banners {
		responseBanners[i] = s.toBannerResponse(banner)
	}

	return responseBanners, nil
}

// Helper methods
func getDisplayOrderOrDefault(order *int, defaultOrder int) int {
	if order != nil {
		return *order
	}
	return defaultOrder
}

func (s *BannerServiceImpl) toBannerResponse(banner *models.Banner) *dto.BannerResponse {
	position := string(banner.Position)
	if position == "" {
		position = "internal"
	}

	return &dto.BannerResponse{
		ID:           banner.ID,
		Title:        banner.Title,
		Image:        cdn.GetFullFileURL(banner.Image), // Construct full URL for frontend
		LinkURL:      banner.LinkURL,
		Alt:          banner.Alt,
		DisplayOrder: banner.DisplayOrder,
		Position:     position,
		Status:       string(banner.Status),
		CreatedAt:    banner.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    banner.UpdatedAt.Format(time.RFC3339),
	}
}