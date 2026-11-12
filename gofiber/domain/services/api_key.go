package services

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
)

type APIKeyService interface {
	// Generate new API key
	GenerateAPIKey(ctx context.Context, userID uuid.UUID, req dto.CreateAPIKeyRequest) (*dto.CreateAPIKeyResponse, error)

	// Validate API key
	ValidateAPIKey(ctx context.Context, apiKey string) (*dto.ValidateAPIKeyResponse, error)

	// List user's API keys
	ListAPIKeys(ctx context.Context, userID uuid.UUID, req dto.ListAPIKeysRequest) (*dto.ListAPIKeysResponse, error)

	// Get API key details
	GetAPIKey(ctx context.Context, userID uuid.UUID, keyID uuid.UUID) (*dto.APIKeyResponse, error)

	// Update API key settings
	UpdateAPIKey(ctx context.Context, userID uuid.UUID, keyID uuid.UUID, req dto.UpdateAPIKeyRequest) (*dto.APIKeyResponse, error)

	// Revoke (soft delete) API key
	RevokeAPIKey(ctx context.Context, userID uuid.UUID, keyID uuid.UUID) error

	// Record API key usage
	RecordUsage(ctx context.Context, keyID uuid.UUID) error
}