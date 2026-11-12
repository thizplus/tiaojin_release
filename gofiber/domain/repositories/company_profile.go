package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type CompanyProfileRepository interface {
	Create(ctx context.Context, entity *models.CompanyProfile) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.CompanyProfile, error)
	Update(ctx context.Context, id uuid.UUID, entity *models.CompanyProfile) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.CompanyProfile, error)
	Count(ctx context.Context) (int64, error)

	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.CompanyProfile, error)
	GetByTaxID(ctx context.Context, taxID string) (*models.CompanyProfile, error)
	SearchByName(ctx context.Context, name string, offset, limit int) ([]*models.CompanyProfile, error)
	UpdateLogo(ctx context.Context, id uuid.UUID, logoURL string) error

	// Methods for admin approval workflow
	ListPendingProfiles(ctx context.Context, offset, limit int) ([]*models.CompanyProfile, error)
	CountPendingProfiles(ctx context.Context) (int64, error)

	// Public
	GetRandomLogos(ctx context.Context, limit int) ([]string, error)
}