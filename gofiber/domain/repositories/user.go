package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, entity *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, id uuid.UUID, entity *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.User, error)
	Count(ctx context.Context) (int64, error)

	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByUsernameWithProfile(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByEmailWithProfile(ctx context.Context, email string) (*models.User, error)
	SearchByUsername(ctx context.Context, username string, offset, limit int) ([]*models.User, error)
	SearchByEmail(ctx context.Context, email string, offset, limit int) ([]*models.User, error)
	GetActiveUsers(ctx context.Context, offset, limit int) ([]*models.User, error)
	CountByStatus(ctx context.Context, status models.UserStatus) (int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.UserStatus) error
}