package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// APIKey represents an API key for external access
type APIKey struct {
	ID          uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	User        *User           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	KeyHash     string          `gorm:"uniqueIndex;not null;size:128" json:"-"` // SHA-256 hash of the key
	Name        string          `gorm:"not null;size:100" json:"name"`
	Description string          `gorm:"size:500" json:"description"`
	IsActive    bool            `gorm:"default:true;not null" json:"is_active"`
	LastUsedAt  *time.Time      `gorm:"index" json:"last_used_at"`
	ExpiresAt   *time.Time      `gorm:"index" json:"expires_at"`
	IPWhitelist pq.StringArray  `gorm:"type:text[]" json:"ip_whitelist"` // Optional IP whitelist
	Permissions pq.StringArray  `gorm:"type:text[]" json:"permissions"`  // e.g., ["tours:read", "cities:read"]
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName sets the table name for APIKey
func (APIKey) TableName() string {
	return "api_keys"
}

// IsExpired checks if the API key has expired
func (a *APIKey) IsExpired() bool {
	if a.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*a.ExpiresAt)
}

// IsValid checks if the API key is valid (active and not expired)
func (a *APIKey) IsValid() bool {
	return a.IsActive && !a.IsExpired()
}

// CanAccessIP checks if the given IP is allowed
func (a *APIKey) CanAccessIP(ip string) bool {
	// If no whitelist, allow all IPs
	if len(a.IPWhitelist) == 0 {
		return true
	}

	// Check if IP is in whitelist
	for _, allowedIP := range a.IPWhitelist {
		if allowedIP == ip {
			return true
		}
	}

	return false
}

// HasPermission checks if the API key has a specific permission
func (a *APIKey) HasPermission(permission string) bool {
	// If no permissions specified, deny all
	if len(a.Permissions) == 0 {
		return false
	}

	for _, p := range a.Permissions {
		if p == permission || p == "*" {
			return true
		}
	}

	return false
}