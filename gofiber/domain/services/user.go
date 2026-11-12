package services

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
)


type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
	GetUserWithProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	UpdateUserStatus(ctx context.Context, userID uuid.UUID, req dto.UpdateUserStatusRequest) (*dto.UserResponse, error)
	ChangeUserRole(ctx context.Context, userID uuid.UUID, req dto.ChangeUserRoleRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	ListUsers(ctx context.Context, req dto.UserListRequest) (*dto.UserListResponse, error)
	GetUserStats(ctx context.Context) (*dto.UserStatsResponse, error)
	ApproveUser(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	SuspendUser(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	ReactivateUser(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
}