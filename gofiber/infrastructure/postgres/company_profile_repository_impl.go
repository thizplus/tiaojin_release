package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type CompanyProfileRepositoryImpl struct {
	db *gorm.DB
}

func NewCompanyProfileRepository(db *gorm.DB) repositories.CompanyProfileRepository {
	return &CompanyProfileRepositoryImpl{db: db}
}

func (r *CompanyProfileRepositoryImpl) Create(ctx context.Context, entity *models.CompanyProfile) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create company profile: %w", err)
	}
	return nil
}

func (r *CompanyProfileRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.CompanyProfile, error) {
	var profile models.CompanyProfile
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&profile).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get company profile by ID: %w", err)
	}

	return &profile, nil
}

func (r *CompanyProfileRepositoryImpl) Update(ctx context.Context, id uuid.UUID, entity *models.CompanyProfile) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(entity)

	if result.Error != nil {
		return fmt.Errorf("failed to update company profile: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("company profile with ID %s not found", id)
	}

	return nil
}

func (r *CompanyProfileRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.CompanyProfile{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete company profile: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("company profile with ID %s not found", id)
	}

	return nil
}

func (r *CompanyProfileRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.CompanyProfile, error) {
	var profiles []*models.CompanyProfile

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&profiles).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list company profiles: %w", err)
	}

	return profiles, nil
}

func (r *CompanyProfileRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.CompanyProfile{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count company profiles: %w", err)
	}

	return count, nil
}

func (r *CompanyProfileRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.CompanyProfile, error) {
	var profile models.CompanyProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get company profile by user ID: %w", err)
	}

	return &profile, nil
}

func (r *CompanyProfileRepositoryImpl) GetByTaxID(ctx context.Context, taxID string) (*models.CompanyProfile, error) {
	var profile models.CompanyProfile
	err := r.db.WithContext(ctx).Where("tax_id = ?", taxID).First(&profile).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get company profile by tax ID: %w", err)
	}

	return &profile, nil
}

func (r *CompanyProfileRepositoryImpl) SearchByName(ctx context.Context, name string, offset, limit int) ([]*models.CompanyProfile, error) {
	var profiles []*models.CompanyProfile
	searchPattern := "%" + name + "%"

	err := r.db.WithContext(ctx).
		Where("name ILIKE ? OR tax_id ILIKE ? OR phone ILIKE ? OR address ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern).
		Order("name ASC").
		Offset(offset).
		Limit(limit).
		Find(&profiles).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search company profiles: %w", err)
	}

	return profiles, nil
}

func (r *CompanyProfileRepositoryImpl) UpdateLogo(ctx context.Context, id uuid.UUID, logoURL string) error {
	result := r.db.WithContext(ctx).
		Model(&models.CompanyProfile{}).
		Where("id = ?", id).
		Update("logo_url", logoURL)

	if result.Error != nil {
		return fmt.Errorf("failed to update company profile logo: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("company profile with ID %s not found", id)
	}

	return nil
}

// ListPendingProfiles lists company profiles for users with pending status
func (r *CompanyProfileRepositoryImpl) ListPendingProfiles(ctx context.Context, offset, limit int) ([]*models.CompanyProfile, error) {
	var profiles []*models.CompanyProfile

	err := r.db.WithContext(ctx).
		Joins("JOIN users ON users.id = company_profiles.user_id").
		Where("users.status = ?", models.UserStatusPending).
		Order("company_profiles.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&profiles).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list pending company profiles: %w", err)
	}

	return profiles, nil
}

// CountPendingProfiles counts company profiles for users with pending status
func (r *CompanyProfileRepositoryImpl) CountPendingProfiles(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.CompanyProfile{}).
		Joins("JOIN users ON users.id = company_profiles.user_id").
		Where("users.status = ?", models.UserStatusPending).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count pending company profiles: %w", err)
	}

	return count, nil
}

// GetRandomLogos ดึง logo URLs แบบ random สำหรับแสดงหน้าแรก
func (r *CompanyProfileRepositoryImpl) GetRandomLogos(ctx context.Context, limit int) ([]string, error) {
	var logos []string

	err := r.db.WithContext(ctx).
		Model(&models.CompanyProfile{}).
		Where("logo IS NOT NULL AND logo != ''").
		Order("RANDOM()").
		Limit(limit).
		Pluck("logo", &logos).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get random logos: %w", err)
	}

	return logos, nil
}

// WithTx returns a new repository instance with the given transaction
func (r *CompanyProfileRepositoryImpl) WithTx(tx *gorm.DB) repositories.CompanyProfileRepository {
	return &CompanyProfileRepositoryImpl{db: tx}
}