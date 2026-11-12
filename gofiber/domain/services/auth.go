package services

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.LoginResponse, error)
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	UpdatePassword(ctx context.Context, req dto.UpdatePasswordRequest, currentUserRole string) error
	ValidateToken(ctx context.Context, token string) (*dto.ValidateTokenResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	RevokeAllTokens(ctx context.Context, userID uuid.UUID) error
}