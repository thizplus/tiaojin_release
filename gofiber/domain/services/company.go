package services

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
)


type CompanyService interface {
	CreateProfile(ctx context.Context, req dto.CreateCompanyProfileRequest) (*dto.CompanyProfileResponse, error)
	GetProfileByID(ctx context.Context, profileID uuid.UUID) (*dto.CompanyProfileResponse, error)
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*dto.CompanyProfileResponse, error)
	GetProfileByTaxID(ctx context.Context, taxID string) (*dto.CompanyProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateCompanyProfileRequest) (*dto.CompanyProfileResponse, error)
	AdminUpdateProfile(ctx context.Context, profileID uuid.UUID, req dto.UpdateCompanyProfileRequest) (*dto.CompanyProfileResponse, error)
	UploadLogo(ctx context.Context, req dto.UploadLogoRequest) (*dto.CompanyProfileResponse, error)
	DeleteLogo(ctx context.Context, userID uuid.UUID) (*dto.CompanyProfileResponse, error)
	UploadPdfHeaderImage(ctx context.Context, req dto.UploadPdfHeaderImageRequest) (*dto.CompanyProfileResponse, error)
	DeletePdfHeaderImage(ctx context.Context, userID uuid.UUID) (*dto.CompanyProfileResponse, error)
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
	ListProfiles(ctx context.Context, req dto.CompanyProfileListRequest) (*dto.CompanyProfileListResponse, error)
	SearchProfiles(ctx context.Context, query string, offset, limit int) (*dto.CompanyProfileListResponse, error)
	ValidateTaxID(ctx context.Context, taxID string, excludeUserID *uuid.UUID) (bool, error)
	ApproveProfile(ctx context.Context, profileID uuid.UUID) (*dto.ProfileApprovalResponse, error)
	GetRandomLogos(ctx context.Context, limit int) ([]string, error)
}