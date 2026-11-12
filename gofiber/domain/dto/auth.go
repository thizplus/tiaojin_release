package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Username       string                         `json:"username" validate:"required,min=3,max=100"`
	Email          *string                        `json:"email,omitempty" validate:"omitempty,email"` // Optional, สำหรับ reset password
	Password       string                         `json:"password" validate:"required,min=6"`
	Role           string                         `json:"role,omitempty" validate:"omitempty,oneof=super admin agent sales"`
	CompanyProfile *RegisterCompanyProfileRequest `json:"company_profile,omitempty"`
}

type RegisterCompanyProfileRequest struct {
	Name     string `json:"name" validate:"required"`
	TaxID    string `json:"tax_id" validate:"required,len=13"`
	Phone    string `json:"phone" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Province string `json:"province" validate:"required"`
	Zipcode  string `json:"zipcode" validate:"required,len=5"`
	ATTA     string `json:"atta" validate:"required"`
	DBD      string `json:"dbd" validate:"required"`
	PP20     string `json:"pp20" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username,omitempty"` // Login ด้วย username
	Email    string `json:"email,omitempty"`    // หรือ Login ด้วย email
	Password string `json:"password" validate:"required"`
}

// GetIdentifier returns username or email for login
func (r *LoginRequest) GetIdentifier() string {
	if r.Username != "" {
		return r.Username
	}
	return r.Email
}

type LoginResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	Email        *string   `json:"email,omitempty"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdatePasswordRequest struct {
	UserID      *uuid.UUID `json:"user_id,omitempty"`                   // Optional: for admin to change other user's password (also accepts "id")
	OldPassword string     `json:"old_password,omitempty"`              // Optional: not required when admin changes other user's password
	NewPassword string     `json:"new_password" validate:"required,min=6"`
}

// Custom UnmarshalJSON to support both "id" and "user_id" fields
func (r *UpdatePasswordRequest) UnmarshalJSON(data []byte) error {
	type Alias UpdatePasswordRequest
	aux := &struct {
		ID *uuid.UUID `json:"id,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// If "id" is provided but "user_id" is not, use "id" as "user_id"
	if aux.ID != nil && r.UserID == nil {
		r.UserID = aux.ID
	}

	return nil
}

type ValidateTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

type ValidateTokenResponse struct {
	Valid    bool      `json:"valid"`
	UserID   uuid.UUID `json:"user_id,omitempty"`
	Username string    `json:"username,omitempty"`
	Email    *string   `json:"email,omitempty"`
	Role     string    `json:"role,omitempty"`
}