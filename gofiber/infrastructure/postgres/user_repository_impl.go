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

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("CompanyProfile").
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetByUsernameWithProfile(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("CompanyProfile").
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get user by username with profile: %w", err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, id uuid.UUID, user *models.User) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(user)

	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %s not found", id)
	}

	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %s not found", id)
	}

	return nil
}

func (r *UserRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
	var users []*models.User

	err := r.db.WithContext(ctx).
		Preload("CompanyProfile").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

func (r *UserRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

func (r *UserRepositoryImpl) GetByEmailWithProfile(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("CompanyProfile").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get user by email with profile: %w", err)
	}

	return &user, nil
}

func (r *UserRepositoryImpl) SearchByUsername(ctx context.Context, username string, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	searchPattern := "%" + username + "%"

	err := r.db.WithContext(ctx).
		Preload("CompanyProfile").
		Where("username ILIKE ?", searchPattern).
		Order("username ASC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search users by username: %w", err)
	}

	return users, nil
}

func (r *UserRepositoryImpl) SearchByEmail(ctx context.Context, email string, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	searchPattern := "%" + email + "%"

	err := r.db.WithContext(ctx).
		Preload("CompanyProfile").
		Where("email ILIKE ?", searchPattern).
		Order("email ASC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search users by email: %w", err)
	}

	return users, nil
}

func (r *UserRepositoryImpl) GetActiveUsers(ctx context.Context, offset, limit int) ([]*models.User, error) {
	var users []*models.User

	err := r.db.WithContext(ctx).
		Preload("CompanyProfile").
		Where("status = ?", models.UserStatusActive).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}

	return users, nil
}

func (r *UserRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status models.UserStatus) error {
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update user status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %s not found", id)
	}

	return nil
}

func (r *UserRepositoryImpl) CountByStatus(ctx context.Context, status models.UserStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("status = ? AND deleted_at IS NULL", status).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count users by status: %w", err)
	}

	return count, nil
}

// WithTx returns a new repository instance with the given transaction
func (r *UserRepositoryImpl) WithTx(tx *gorm.DB) repositories.UserRepository {
	return &UserRepositoryImpl{db: tx}
}