package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CreateAPIKeyRequest represents request to create a new API key
type CreateAPIKeyRequest struct {
	Name        string     `json:"name" validate:"required,min=3,max=100"`
	Description string     `json:"description,omitempty" validate:"max=500"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	IPWhitelist []string   `json:"ip_whitelist,omitempty"`
	Permissions []string   `json:"permissions,omitempty"`
}

// CreateAPIKeyResponse represents response after creating API key
// The actual key is only shown once during creation
type CreateAPIKeyResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	APIKey      string     `json:"api_key"` // Plain text key - shown only once
	IsActive    bool       `json:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at"`
	IPWhitelist []string   `json:"ip_whitelist"`
	Permissions []string   `json:"permissions"`
	CreatedAt   time.Time  `json:"created_at"`
}

// APIKeyResponse represents API key information (without the actual key)
type APIKeyResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	KeyPreview  string     `json:"key_preview"` // e.g., "tiaojin_abc...xyz"
	IsActive    bool       `json:"is_active"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	IPWhitelist []string   `json:"ip_whitelist"`
	Permissions []string   `json:"permissions"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UpdateAPIKeyRequest represents request to update API key settings
type UpdateAPIKeyRequest struct {
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=500"`
	IsActive    *bool      `json:"is_active,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	IPWhitelist []string   `json:"ip_whitelist,omitempty"`
	Permissions []string   `json:"permissions,omitempty"`
}

// UnmarshalJSON custom unmarshaler to properly handle boolean pointers
func (r *UpdateAPIKeyRequest) UnmarshalJSON(data []byte) error {
	// First unmarshal into a map to check if is_active is present
	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	// Unmarshal normally for other fields
	type Alias UpdateAPIKeyRequest
	aux := (*Alias)(r)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle is_active explicitly if present in JSON
	if isActiveVal, exists := rawMap["is_active"]; exists {
		if boolVal, ok := isActiveVal.(bool); ok {
			r.IsActive = &boolVal
		}
	}

	return nil
}

// ListAPIKeysRequest represents request to list API keys
type ListAPIKeysRequest struct {
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
	IsActive   *bool  `json:"is_active,omitempty"`
	Search     string `json:"search,omitempty"`
	ShowExpired bool  `json:"show_expired"`
}

// ListAPIKeysResponse represents response with list of API keys
type ListAPIKeysResponse struct {
	APIKeys    []*APIKeyResponse `json:"api_keys"`
	Pagination PaginationMeta    `json:"pagination"`
}

// ValidateAPIKeyRequest represents request to validate an API key
type ValidateAPIKeyRequest struct {
	APIKey string `json:"api_key" validate:"required"`
}

// ValidateAPIKeyResponse represents validation result
type ValidateAPIKeyResponse struct {
	Valid       bool       `json:"valid"`
	UserID      uuid.UUID  `json:"user_id,omitempty"`
	KeyID       uuid.UUID  `json:"key_id,omitempty"`
	Permissions []string   `json:"permissions,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}