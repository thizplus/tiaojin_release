package serviceimpl

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
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


const (
	MaxLogoSize = 5 * 1024 * 1024 // 5MB
)

type CompanyServiceImpl struct {
	companyRepo repositories.CompanyProfileRepository
	userRepo    repositories.UserRepository
}

func NewCompanyService(
	companyRepo repositories.CompanyProfileRepository,
	userRepo repositories.UserRepository,
) services.CompanyService {
	return &CompanyServiceImpl{
		companyRepo: companyRepo,
		userRepo:    userRepo,
	}
}

func (s *CompanyServiceImpl) CreateProfile(ctx context.Context, req dto.CreateCompanyProfileRequest) (*dto.CompanyProfileResponse, error) {
	// Business validation - check if user exists
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return nil, fmt.Errorf("failed to validate user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Business validation - check if company profile already exists
	if existing, err := s.companyRepo.GetByUserID(ctx, req.UserID); err != nil {
		log.Printf("Error checking existing company profile: %v", err)
		return nil, fmt.Errorf("failed to validate company profile: %w", err)
	} else if existing != nil {
		return nil, apperrors.ErrCompanyProfileExists
	}

	// Business validation - check if tax ID already exists
	if existing, err := s.companyRepo.GetByTaxID(ctx, req.TaxID); err != nil {
		log.Printf("Error checking tax ID existence: %v", err)
		return nil, fmt.Errorf("failed to validate tax ID: %w", err)
	} else if existing != nil {
		return nil, apperrors.ErrTaxIDAlreadyExists
	}

	// Business validation - validate tax ID format (13 digits)
	if !isValidTaxID(req.TaxID) {
		return nil, apperrors.ErrInvalidTaxIDFormat
	}

	// Business validation - validate zipcode format (5 digits)
	if !isValidZipcode(req.Zipcode) {
		return nil, apperrors.ErrInvalidZipcodeFormat
	}

	// Create domain model
	profile := &models.CompanyProfile{
		UserID:    req.UserID,
		Name:      req.Name,
		TaxID:     req.TaxID,
		Phone:     req.Phone,
		Address:   req.Address,
		Province:  req.Province,
		Zipcode:   req.Zipcode,
		ATTA:      req.ATTA,
		DBD:       req.DBD,
		PP20:      req.PP20,
		Logo:      req.Logo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create company profile
	if err := s.companyRepo.Create(ctx, profile); err != nil {
		log.Printf("Error creating company profile: %v", err)
		return nil, fmt.Errorf("failed to create company profile: %w", err)
	}

	return s.toCompanyProfileResponse(profile), nil
}

func (s *CompanyServiceImpl) GetProfileByID(ctx context.Context, profileID uuid.UUID) (*dto.CompanyProfileResponse, error) {
	profile, err := s.companyRepo.GetByID(ctx, profileID)
	if err != nil {
		log.Printf("Error retrieving company profile: %v", err)
		return nil, fmt.Errorf("failed to get company profile: %w", err)
	}

	if profile == nil {
		return nil, apperrors.ErrCompanyProfileNotFound
	}

	return s.toCompanyProfileResponse(profile), nil
}

func (s *CompanyServiceImpl) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*dto.CompanyProfileResponse, error) {
	profile, err := s.companyRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error retrieving company profile by user ID: %v", err)
		return nil, fmt.Errorf("failed to get company profile: %w", err)
	}

	if profile == nil {
		return nil, apperrors.ErrCompanyProfileNotFound
	}

	return s.toCompanyProfileResponse(profile), nil
}

func (s *CompanyServiceImpl) GetProfileByTaxID(ctx context.Context, taxID string) (*dto.CompanyProfileResponse, error) {
	// Business validation - validate tax ID format
	if !isValidTaxID(taxID) {
		return nil, apperrors.ErrInvalidTaxIDFormat
	}

	profile, err := s.companyRepo.GetByTaxID(ctx, taxID)
	if err != nil {
		log.Printf("Error retrieving company profile by tax ID: %v", err)
		return nil, fmt.Errorf("failed to get company profile: %w", err)
	}

	if profile == nil {
		return nil, apperrors.ErrCompanyProfileNotFound
	}

	return s.toCompanyProfileResponse(profile), nil
}

func (s *CompanyServiceImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateCompanyProfileRequest) (*dto.CompanyProfileResponse, error) {
	// Get existing profile
	profile, err := s.companyRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error retrieving company profile: %v", err)
		return nil, fmt.Errorf("failed to get company profile: %w", err)
	}

	if profile == nil {
		return nil, apperrors.ErrCompanyProfileNotFound
	}

	// Business validation - if tax ID is being updated, check uniqueness
	if req.TaxID != nil && *req.TaxID != profile.TaxID {
		if !isValidTaxID(*req.TaxID) {
			return nil, apperrors.ErrInvalidTaxIDFormat
		}

		if existing, err := s.companyRepo.GetByTaxID(ctx, *req.TaxID); err != nil {
			log.Printf("Error checking tax ID existence: %v", err)
			return nil, fmt.Errorf("failed to validate tax ID: %w", err)
		} else if existing != nil {
			return nil, apperrors.ErrTaxIDAlreadyExists
		}
	}

	// Business validation - validate zipcode format
	if req.Zipcode != nil && !isValidZipcode(*req.Zipcode) {
		return nil, apperrors.ErrInvalidZipcodeFormat
	}

	// Update fields if provided
	if req.Name != nil {
		profile.Name = *req.Name
	}
	if req.TaxID != nil {
		profile.TaxID = *req.TaxID
	}
	if req.Phone != nil {
		profile.Phone = *req.Phone
	}
	if req.Address != nil {
		profile.Address = *req.Address
	}
	if req.Province != nil {
		profile.Province = *req.Province
	}
	if req.Zipcode != nil {
		profile.Zipcode = *req.Zipcode
	}
	if req.ATTA != nil {
		profile.ATTA = *req.ATTA
	}
	if req.DBD != nil {
		profile.DBD = *req.DBD
	}
	if req.PP20 != nil {
		profile.PP20 = *req.PP20
	}
	if req.Logo != nil {
		profile.Logo = req.Logo
	}

	profile.UpdatedAt = time.Now()

	// Update profile
	if err := s.companyRepo.Update(ctx, profile.ID, profile); err != nil {
		log.Printf("Error updating company profile: %v", err)
		return nil, fmt.Errorf("failed to update company profile: %w", err)
	}

	return s.toCompanyProfileResponse(profile), nil
}

// AdminUpdateProfile — admin แก้ไข profile โดยใช้ profileID (ไม่เช็ค ownership)
func (s *CompanyServiceImpl) AdminUpdateProfile(ctx context.Context, profileID uuid.UUID, req dto.UpdateCompanyProfileRequest) (*dto.CompanyProfileResponse, error) {
	profile, err := s.companyRepo.GetByID(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company profile: %w", err)
	}
	if profile == nil {
		return nil, apperrors.ErrCompanyProfileNotFound
	}

	// Business validation - tax ID uniqueness
	if req.TaxID != nil && *req.TaxID != profile.TaxID {
		if !isValidTaxID(*req.TaxID) {
			return nil, apperrors.ErrInvalidTaxIDFormat
		}
		if existing, err := s.companyRepo.GetByTaxID(ctx, *req.TaxID); err != nil {
			return nil, fmt.Errorf("failed to validate tax ID: %w", err)
		} else if existing != nil && existing.ID != profileID {
			return nil, apperrors.ErrTaxIDAlreadyExists
		}
	}

	if req.Zipcode != nil && !isValidZipcode(*req.Zipcode) {
		return nil, apperrors.ErrInvalidZipcodeFormat
	}

	// Update fields
	if req.Name != nil {
		profile.Name = *req.Name
	}
	if req.TaxID != nil {
		profile.TaxID = *req.TaxID
	}
	if req.Phone != nil {
		profile.Phone = *req.Phone
	}
	if req.Address != nil {
		profile.Address = *req.Address
	}
	if req.Province != nil {
		profile.Province = *req.Province
	}
	if req.Zipcode != nil {
		profile.Zipcode = *req.Zipcode
	}
	if req.Logo != nil {
		profile.Logo = req.Logo
	}

	profile.UpdatedAt = time.Now()

	if err := s.companyRepo.Update(ctx, profile.ID, profile); err != nil {
		return nil, fmt.Errorf("failed to update company profile: %w", err)
	}

	return s.toCompanyProfileResponse(profile), nil
}

func (s *CompanyServiceImpl) UploadLogo(ctx context.Context, req dto.UploadLogoRequest) (*dto.CompanyProfileResponse, error) {
	// Business validation - check file size
	if len(req.LogoData) > MaxLogoSize {
		return nil, apperrors.ErrLogoTooLarge
	}

	// Business validation - check file type
	if !isValidImageType(req.FileType) {
		return nil, apperrors.ErrInvalidLogoFormat
	}

	// Get existing profile
	profile, err := s.companyRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		log.Printf("Error retrieving company profile: %v", err)
		return nil, fmt.Errorf("failed to get company profile: %w", err)
	}

	if profile == nil {
		return nil, apperrors.ErrCompanyProfileNotFound
	}

	// Generate logo URL/path (in real implementation, you would upload to storage service)
	logoURL := s.generateLogoPath(req.UserID, req.FileName)

	// Update logo URL
	if err := s.companyRepo.UpdateLogo(ctx, profile.ID, logoURL); err != nil {
		log.Printf("Error updating company profile logo: %v", err)
		return nil, fmt.Errorf("failed to update logo: %w", err)
	}

	// Update local profile object
	profile.Logo = &logoURL
	profile.UpdatedAt = time.Now()

	return s.toCompanyProfileResponse(profile), nil
}

func (s *CompanyServiceImpl) DeleteLogo(ctx context.Context, userID uuid.UUID) (*dto.CompanyProfileResponse, error) {
	// Get existing profile
	profile, err := s.companyRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error retrieving company profile: %v", err)
		return nil, fmt.Errorf("failed to get company profile: %w", err)
	}

	if profile == nil {
		return nil, apperrors.ErrCompanyProfileNotFound
	}

	// Remove logo
	if err := s.companyRepo.UpdateLogo(ctx, profile.ID, ""); err != nil {
		log.Printf("Error removing company profile logo: %v", err)
		return nil, fmt.Errorf("failed to remove logo: %w", err)
	}

	// Update local profile object
	profile.Logo = nil
	profile.UpdatedAt = time.Now()

	return s.toCompanyProfileResponse(profile), nil
}

func (s *CompanyServiceImpl) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	// Get existing profile
	profile, err := s.companyRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error retrieving company profile: %v", err)
		return fmt.Errorf("failed to get company profile: %w", err)
	}

	if profile == nil {
		return apperrors.ErrCompanyProfileNotFound
	}

	// Delete profile
	if err := s.companyRepo.Delete(ctx, profile.ID); err != nil {
		log.Printf("Error deleting company profile: %v", err)
		return fmt.Errorf("failed to delete company profile: %w", err)
	}

	return nil
}

func (s *CompanyServiceImpl) ListProfiles(ctx context.Context, req dto.CompanyProfileListRequest) (*dto.CompanyProfileListResponse, error) {
	var profiles []*models.CompanyProfile
	var total int64
	var err error

	// If PendingOnly is true, filter by users with pending status
	if req.PendingOnly {
		profiles, err = s.companyRepo.ListPendingProfiles(ctx, req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error listing pending company profiles: %v", err)
			return nil, fmt.Errorf("failed to list pending company profiles: %w", err)
		}

		total, err = s.companyRepo.CountPendingProfiles(ctx)
		if err != nil {
			log.Printf("Error counting pending company profiles: %v", err)
			return nil, fmt.Errorf("failed to count pending company profiles: %w", err)
		}
	} else if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		// If search query is provided, use search method
		profiles, err = s.companyRepo.SearchByName(ctx, strings.TrimSpace(*req.Search), req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error searching company profiles: %v", err)
			return nil, fmt.Errorf("failed to search company profiles: %w", err)
		}

		// For search, we need to count separately (simplified approach)
		total = int64(len(profiles))
	} else {
		// List all profiles
		profiles, err = s.companyRepo.List(ctx, req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error listing company profiles: %v", err)
			return nil, fmt.Errorf("failed to list company profiles: %w", err)
		}

		// Get total count
		total, err = s.companyRepo.Count(ctx)
		if err != nil {
			log.Printf("Error counting company profiles: %v", err)
			return nil, fmt.Errorf("failed to count company profiles: %w", err)
		}
	}

	// Convert to response DTOs
	responseProfiles := make([]*dto.CompanyProfileResponse, len(profiles))
	for i, profile := range profiles {
		responseProfiles[i] = s.toCompanyProfileResponse(profile)
	}

	return &dto.CompanyProfileListResponse{
		Profiles:   responseProfiles,
		Pagination: utils.CalculatePagination(total, req.Offset, req.Limit),
	}, nil
}

func (s *CompanyServiceImpl) SearchProfiles(ctx context.Context, query string, offset, limit int) (*dto.CompanyProfileListResponse, error) {
	// Business validation - clean search query
	query = strings.TrimSpace(query)
	if query == "" {
		return &dto.CompanyProfileListResponse{
			Profiles:   []*dto.CompanyProfileResponse{},
			Pagination: utils.CalculatePagination(0, offset, limit),
		}, nil
	}

	profiles, err := s.companyRepo.SearchByName(ctx, query, offset, limit)
	if err != nil {
		log.Printf("Error searching company profiles: %v", err)
		return nil, fmt.Errorf("failed to search company profiles: %w", err)
	}

	// Convert to response DTOs
	responseProfiles := make([]*dto.CompanyProfileResponse, len(profiles))
	for i, profile := range profiles {
		responseProfiles[i] = s.toCompanyProfileResponse(profile)
	}

	// For search results, total is approximate
	total := int64(len(profiles))

	return &dto.CompanyProfileListResponse{
		Profiles:   responseProfiles,
		Pagination: utils.CalculatePagination(total, offset, limit),
	}, nil
}

func (s *CompanyServiceImpl) ValidateTaxID(ctx context.Context, taxID string, excludeUserID *uuid.UUID) (bool, error) {
	// Business validation - validate format first
	if !isValidTaxID(taxID) {
		return false, nil
	}

	// Check if tax ID exists
	existing, err := s.companyRepo.GetByTaxID(ctx, taxID)
	if err != nil {
		log.Printf("Error checking tax ID existence: %v", err)
		return false, fmt.Errorf("failed to validate tax ID: %w", err)
	}

	// If no existing profile found, tax ID is available
	if existing == nil {
		return true, nil
	}

	// If excludeUserID is provided and matches, tax ID is available for that user
	if excludeUserID != nil && existing.UserID == *excludeUserID {
		return true, nil
	}

	// Tax ID already exists for another user
	return false, nil
}

// Helper methods
func (s *CompanyServiceImpl) toCompanyProfileResponse(profile *models.CompanyProfile) *dto.CompanyProfileResponse {
	return &dto.CompanyProfileResponse{
		ID:        profile.ID,
		UserID:    profile.UserID,
		Name:      profile.Name,
		TaxID:     profile.TaxID,
		Phone:     profile.Phone,
		Address:   profile.Address,
		Province:  profile.Province,
		Zipcode:   profile.Zipcode,
		ATTA:      cdn.GetFullFileURL(profile.ATTA),      // Construct full URL for frontend
		DBD:       cdn.GetFullFileURL(profile.DBD),       // Construct full URL for frontend
		PP20:      cdn.GetFullFileURL(profile.PP20),      // Construct full URL for frontend
		Logo:      cdn.GetFullFileURLPtr(profile.Logo),   // Construct full URL for frontend
		CreatedAt: profile.CreatedAt.Format(time.RFC3339),
		UpdatedAt: profile.UpdatedAt.Format(time.RFC3339),
	}
}

func isValidTaxID(taxID string) bool {
	// Tax ID should be exactly 13 digits
	if len(taxID) != 13 {
		return false
	}

	for _, char := range taxID {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

func isValidZipcode(zipcode string) bool {
	// Zipcode should be exactly 5 digits
	if len(zipcode) != 5 {
		return false
	}

	for _, char := range zipcode {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

func isValidImageType(fileType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/png",
		"image/webp",
	}

	for _, validType := range validTypes {
		if fileType == validType {
			return true
		}
	}

	return false
}

func (s *CompanyServiceImpl) generateLogoPath(userID uuid.UUID, fileName string) string {
	// Generate a unique file path for the logo
	ext := filepath.Ext(fileName)
	cleanExt := strings.ToLower(ext)

	// Validate extension against MIME type
	if cleanExt == "" {
		cleanExt = ".jpg" // default
	}

	timestamp := time.Now().Unix()
	return fmt.Sprintf("logos/%s/%d%s", userID.String(), timestamp, cleanExt)
}

func (s *CompanyServiceImpl) ApproveProfile(ctx context.Context, profileID uuid.UUID) (*dto.ProfileApprovalResponse, error) {
	// Get existing profile
	profile, err := s.companyRepo.GetByID(ctx, profileID)
	if err != nil {
		log.Printf("Error retrieving company profile: %v", err)
		return nil, fmt.Errorf("failed to get company profile: %w", err)
	}

	if profile == nil {
		return nil, apperrors.ErrCompanyProfileNotFound
	}

	// Get user to check current status
	user, err := s.userRepo.GetByID(ctx, profile.UserID)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is already active
	if user.Status == models.UserStatusActive {
		return nil, fmt.Errorf("user is already active")
	}

	// Update user status to active
	if err := s.userRepo.UpdateStatus(ctx, profile.UserID, models.UserStatusActive); err != nil {
		log.Printf("Error updating user status: %v", err)
		return nil, fmt.Errorf("failed to activate user: %w", err)
	}

	return &dto.ProfileApprovalResponse{
		ID:        profile.ID,
		UserID:    profile.UserID,
		Status:    "approved",
		Message:   "Company profile approved and user activated successfully",
		UpdatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func (s *CompanyServiceImpl) GetRandomLogos(ctx context.Context, limit int) ([]string, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	paths, err := s.companyRepo.GetRandomLogos(ctx, limit)
	if err != nil {
		return nil, err
	}

	// แปลง CDN path → full URL (เหมือนที่ทำใน toCompanyProfileResponse)
	urls := make([]string, 0, len(paths))
	for _, path := range paths {
		fullURL := cdn.GetFullFileURL(path)
		if fullURL != "" {
			urls = append(urls, fullURL)
		}
	}
	return urls, nil
}