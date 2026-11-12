package serviceimpl

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
)

const (
	APIKeyPrefix = "tiaojin_"
	APIKeyLength = 40 // length of random part
)

type APIKeyServiceImpl struct {
	apiKeyRepo repositories.APIKeyRepository
	userRepo   repositories.UserRepository
}

func NewAPIKeyService(apiKeyRepo repositories.APIKeyRepository, userRepo repositories.UserRepository) services.APIKeyService {
	return &APIKeyServiceImpl{
		apiKeyRepo: apiKeyRepo,
		userRepo:   userRepo,
	}
}

// GenerateAPIKey creates a new API key for a user
func (s *APIKeyServiceImpl) GenerateAPIKey(ctx context.Context, userID uuid.UUID, req dto.CreateAPIKeyRequest) (*dto.CreateAPIKeyResponse, error) {
	log.Printf("🔑 Generating API key for user %s with name: %s", userID, req.Name)

	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate random API key
	apiKey, err := generateRandomAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	// Hash the API key for storage
	keyHash := hashAPIKey(apiKey)

	// Set default permissions if not provided
	permissions := req.Permissions
	if len(permissions) == 0 {
		// Default permissions for agent: read-only access to tours and cities
		permissions = []string{"tours:read", "cities:read", "periods:read"}
	}

	// Convert to pq.StringArray
	ipWhitelist := pq.StringArray(req.IPWhitelist)
	permissionsArray := pq.StringArray(permissions)

	// Create API key model
	apiKeyModel := &models.APIKey{
		UserID:      userID,
		KeyHash:     keyHash,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		ExpiresAt:   req.ExpiresAt,
		IPWhitelist: ipWhitelist,
		Permissions: permissionsArray,
	}

	// Save to database
	if err := s.apiKeyRepo.Create(ctx, apiKeyModel); err != nil {
		return nil, fmt.Errorf("failed to save API key: %w", err)
	}

	log.Printf("✅ API key created successfully: %s (ID: %s)", req.Name, apiKeyModel.ID)

	// Return response with plain text key (only shown once)
	return &dto.CreateAPIKeyResponse{
		ID:          apiKeyModel.ID,
		Name:        apiKeyModel.Name,
		Description: apiKeyModel.Description,
		APIKey:      apiKey, // Plain text - only shown once!
		IsActive:    apiKeyModel.IsActive,
		ExpiresAt:   apiKeyModel.ExpiresAt,
		IPWhitelist: []string(apiKeyModel.IPWhitelist),
		Permissions: []string(apiKeyModel.Permissions),
		CreatedAt:   apiKeyModel.CreatedAt,
	}, nil
}

// ValidateAPIKey validates an API key and returns user info
func (s *APIKeyServiceImpl) ValidateAPIKey(ctx context.Context, apiKey string) (*dto.ValidateAPIKeyResponse, error) {
	// Check prefix
	if !strings.HasPrefix(apiKey, APIKeyPrefix) {
		return &dto.ValidateAPIKeyResponse{Valid: false}, nil
	}

	// Hash the provided key
	keyHash := hashAPIKey(apiKey)

	// Find API key in database
	apiKeyModel, err := s.apiKeyRepo.GetByKeyHash(ctx, keyHash)
	if err != nil {
		return nil, fmt.Errorf("failed to find API key: %w", err)
	}

	if apiKeyModel == nil {
		return &dto.ValidateAPIKeyResponse{Valid: false}, nil
	}

	// Check if key is valid (active and not expired)
	if !apiKeyModel.IsValid() {
		return &dto.ValidateAPIKeyResponse{Valid: false}, nil
	}

	// Update last used timestamp (async - don't block validation)
	go func() {
		if err := s.apiKeyRepo.UpdateLastUsed(context.Background(), apiKeyModel.ID); err != nil {
			log.Printf("⚠️ Failed to update last used timestamp: %v", err)
		}
	}()

	return &dto.ValidateAPIKeyResponse{
		Valid:       true,
		UserID:      apiKeyModel.UserID,
		KeyID:       apiKeyModel.ID,
		Permissions: []string(apiKeyModel.Permissions),
		ExpiresAt:   apiKeyModel.ExpiresAt,
	}, nil
}

// ListAPIKeys lists all API keys for a user
func (s *APIKeyServiceImpl) ListAPIKeys(ctx context.Context, userID uuid.UUID, req dto.ListAPIKeysRequest) (*dto.ListAPIKeysResponse, error) {
	log.Printf("📋 Listing API keys for user %s", userID)

	// Set default pagination
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Get API keys with filters
	apiKeys, err := s.apiKeyRepo.ListWithFilters(ctx, req, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}

	// Get total count
	totalCount, err := s.apiKeyRepo.CountWithFilters(ctx, req, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count API keys: %w", err)
	}

	// Convert to response DTOs
	apiKeyResponses := make([]*dto.APIKeyResponse, len(apiKeys))
	for i, apiKey := range apiKeys {
		apiKeyResponses[i] = toAPIKeyResponse(apiKey)
	}

	return &dto.ListAPIKeysResponse{
		APIKeys: apiKeyResponses,
		Pagination: dto.PaginationMeta{
			Total:  totalCount,
			Offset: req.Offset,
			Limit:  req.Limit,
		},
	}, nil
}

// GetAPIKey gets details of a specific API key
func (s *APIKeyServiceImpl) GetAPIKey(ctx context.Context, userID uuid.UUID, keyID uuid.UUID) (*dto.APIKeyResponse, error) {
	log.Printf("🔍 Getting API key %s for user %s", keyID, userID)

	apiKey, err := s.apiKeyRepo.GetByID(ctx, keyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	if apiKey == nil {
		return nil, fmt.Errorf("API key not found")
	}

	// Verify ownership
	if apiKey.UserID != userID {
		return nil, fmt.Errorf("access denied: API key belongs to different user")
	}

	return toAPIKeyResponse(apiKey), nil
}

// UpdateAPIKey updates API key settings
func (s *APIKeyServiceImpl) UpdateAPIKey(ctx context.Context, userID uuid.UUID, keyID uuid.UUID, req dto.UpdateAPIKeyRequest) (*dto.APIKeyResponse, error) {
	log.Printf("✏️ Updating API key %s for user %s", keyID, userID)

	// Get existing API key
	apiKey, err := s.apiKeyRepo.GetByID(ctx, keyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	if apiKey == nil {
		return nil, fmt.Errorf("API key not found")
	}

	// Verify ownership
	if apiKey.UserID != userID {
		return nil, fmt.Errorf("access denied: API key belongs to different user")
	}

	// Update fields conditionally
	if req.Name != nil {
		apiKey.Name = *req.Name
	}
	if req.Description != nil {
		apiKey.Description = *req.Description
	}
	if req.IsActive != nil {
		log.Printf("🔄 Updating is_active from %v to %v", apiKey.IsActive, *req.IsActive)
		apiKey.IsActive = *req.IsActive
	}
	if req.ExpiresAt != nil {
		apiKey.ExpiresAt = req.ExpiresAt
	}
	if req.IPWhitelist != nil {
		apiKey.IPWhitelist = pq.StringArray(req.IPWhitelist)
	}
	if req.Permissions != nil {
		apiKey.Permissions = pq.StringArray(req.Permissions)
	}

	apiKey.UpdatedAt = time.Now()

	// Save updates - use Select to force update is_active even if false
	if err := s.apiKeyRepo.Update(ctx, keyID, apiKey); err != nil {
		return nil, fmt.Errorf("failed to update API key: %w", err)
	}

	log.Printf("✅ API key updated successfully: %s", keyID)

	return toAPIKeyResponse(apiKey), nil
}

// RevokeAPIKey revokes (soft deletes) an API key
func (s *APIKeyServiceImpl) RevokeAPIKey(ctx context.Context, userID uuid.UUID, keyID uuid.UUID) error {
	log.Printf("🗑️ Revoking API key %s for user %s", keyID, userID)

	// Get API key to verify ownership
	apiKey, err := s.apiKeyRepo.GetByID(ctx, keyID)
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	if apiKey == nil {
		return fmt.Errorf("API key not found")
	}

	// Verify ownership
	if apiKey.UserID != userID {
		return fmt.Errorf("access denied: API key belongs to different user")
	}

	// Soft delete
	if err := s.apiKeyRepo.SoftDelete(ctx, keyID); err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	log.Printf("✅ API key revoked successfully: %s", keyID)

	return nil
}

// RecordUsage records API key usage (called by middleware)
func (s *APIKeyServiceImpl) RecordUsage(ctx context.Context, keyID uuid.UUID) error {
	return s.apiKeyRepo.UpdateLastUsed(ctx, keyID)
}

// Helper functions

// generateRandomAPIKey generates a random API key
func generateRandomAPIKey() (string, error) {
	bytes := make([]byte, APIKeyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	randomPart := hex.EncodeToString(bytes)[:APIKeyLength]
	return APIKeyPrefix + randomPart, nil
}

// hashAPIKey creates a SHA-256 hash of the API key
func hashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(hash[:])
}

// toAPIKeyResponse converts model to response DTO
func toAPIKeyResponse(apiKey *models.APIKey) *dto.APIKeyResponse {
	// Create key preview (first 10 + last 5 characters)
	keyPreview := APIKeyPrefix + "••••••••••••••••••••••••"
	if len(apiKey.KeyHash) >= 15 {
		keyPreview = APIKeyPrefix + "..." + apiKey.KeyHash[len(apiKey.KeyHash)-5:]
	}

	return &dto.APIKeyResponse{
		ID:          apiKey.ID,
		Name:        apiKey.Name,
		Description: apiKey.Description,
		KeyPreview:  keyPreview,
		IsActive:    apiKey.IsActive,
		LastUsedAt:  apiKey.LastUsedAt,
		ExpiresAt:   apiKey.ExpiresAt,
		IPWhitelist: []string(apiKey.IPWhitelist),
		Permissions: []string(apiKey.Permissions),
		CreatedAt:   apiKey.CreatedAt,
		UpdatedAt:   apiKey.UpdatedAt,
	}
}