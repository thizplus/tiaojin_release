package serviceimpl

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	apperrors "gofiber-template/pkg/errors"
	"gofiber-template/pkg/utils"
)

type AuthServiceImpl struct {
	userRepo        repositories.UserRepository
	companyRepo     repositories.CompanyProfileRepository
	jwtSecret       string
	tokenExpiration time.Duration
}

func NewAuthService(
	userRepo repositories.UserRepository,
	companyRepo repositories.CompanyProfileRepository,
	jwtSecret string,
) services.AuthService {
	return &AuthServiceImpl{
		userRepo:        userRepo,
		companyRepo:     companyRepo,
		jwtSecret:       jwtSecret,
		tokenExpiration: time.Hour * 24 * 7, // 7 days
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error) {
	// Business validation - check if username already exists
	if existing, err := s.userRepo.GetByUsername(ctx, req.Username); err != nil {
		log.Printf("Error checking username existence: %v", err)
		return nil, fmt.Errorf("failed to validate username: %w", err)
	} else if existing != nil {
		return nil, apperrors.ErrUsernameExists
	}

	// Check email uniqueness if provided
	if req.Email != nil && *req.Email != "" {
		if existing, err := s.userRepo.GetByEmail(ctx, *req.Email); err != nil {
			log.Printf("Error checking email existence: %v", err)
			return nil, fmt.Errorf("failed to validate email: %w", err)
		} else if existing != nil {
			return nil, apperrors.ErrEmailExists
		}
	}

	// Business logic - hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, fmt.Errorf("failed to process password: %w", err)
	}

	// Determine user role
	userRole := models.UserRoleAgent // Default role for new registrations
	if req.Role != "" {
		userRole = models.UserRole(req.Role)
	}

	// Create domain models
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      userRole,
		Status:    models.UserStatusPending, // Requires approval
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create user
	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create company profile only if provided
	if req.CompanyProfile != nil {
		companyProfile := &models.CompanyProfile{
			UserID:    user.ID,
			Name:      req.CompanyProfile.Name,
			TaxID:     req.CompanyProfile.TaxID,
			Phone:     req.CompanyProfile.Phone,
			Address:   req.CompanyProfile.Address,
			Province:  req.CompanyProfile.Province,
			Zipcode:   req.CompanyProfile.Zipcode,
			ATTA:      req.CompanyProfile.ATTA,
			DBD:       req.CompanyProfile.DBD,
			PP20:      req.CompanyProfile.PP20,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.companyRepo.Create(ctx, companyProfile); err != nil {
			log.Printf("Error creating company profile: %v", err)
			// Rollback user creation would require transaction
			return nil, fmt.Errorf("failed to create company profile: %w", err)
		}
	}

	// Generate tokens for immediate login
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &dto.LoginResponse{
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Role:         string(user.Role),
		Status:       string(user.Status),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.tokenExpiration.Seconds()),
	}, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get identifier (username or email)
	identifier := req.GetIdentifier()
	if identifier == "" {
		return nil, apperrors.ErrInvalidCredentials
	}

	// Try to find user by username first, then by email
	var user *models.User
	var err error

	user, err = s.userRepo.GetByUsername(ctx, identifier)
	if err != nil || user == nil {
		// Try by email if username not found
		user, err = s.userRepo.GetByEmail(ctx, identifier)
		if err != nil {
			log.Printf("Error retrieving user: %v", err)
			return nil, apperrors.ErrInvalidCredentials
		}
	}

	if user == nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	// Business rule - check if account is not inactive
	if user.Status == models.UserStatusInactive {
		return nil, apperrors.ErrAccountInactive
	}

	// Business logic - verify password
	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &dto.LoginResponse{
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Role:         string(user.Role),
		Status:       string(user.Status),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.tokenExpiration.Seconds()),
	}, nil
}

func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	// Parse and validate refresh token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.ErrInvalidToken
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, apperrors.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, apperrors.ErrInvalidToken
	}

	// Check if token type is refresh
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, apperrors.ErrInvalidToken
	}

	// Get user ID from claims
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, apperrors.ErrInvalidToken
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, apperrors.ErrInvalidToken
	}

	// Get user from database
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, apperrors.ErrUserNotFound
	}

	// Check if user is not inactive
	if user.Status == models.UserStatusInactive {
		return nil, apperrors.ErrAccountInactive
	}

	// Generate new tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &dto.LoginResponse{
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Role:         string(user.Role),
		Status:       string(user.Status),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.tokenExpiration.Seconds()),
	}, nil
}

func (s *AuthServiceImpl) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	// Check if user exists
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("Error retrieving user for password reset: %v", err)
		return fmt.Errorf("failed to process password reset: %w", err)
	}

	if user == nil {
		// Don't reveal if email exists - return success anyway for security
		return nil
	}

	// TODO: Generate password reset token and send email
	// For now, this is a placeholder implementation
	log.Printf("Password reset requested for user: %s", req.Email)

	return nil
}

func (s *AuthServiceImpl) UpdatePassword(ctx context.Context, req dto.UpdatePasswordRequest, currentUserRole string) error {
	// Validate user_id is provided
	if req.UserID == nil {
		return fmt.Errorf("user_id is required")
	}

	// Get target user
	user, err := s.userRepo.GetByID(ctx, *req.UserID)
	if err != nil || user == nil {
		return apperrors.ErrUserNotFound
	}

	// Check if admin/super is changing another user's password
	isAdminChangingOther := (currentUserRole == "admin" || currentUserRole == "super") && req.OldPassword == ""

	// Verify old password (only if not admin changing other user's password)
	if !isAdminChangingOther {
		if req.OldPassword == "" {
			return fmt.Errorf("old_password is required when changing your own password")
		}
		if err := utils.ComparePassword(user.Password, req.OldPassword); err != nil {
			return apperrors.ErrInvalidCredentials
		}
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("Error hashing new password: %v", err)
		return fmt.Errorf("failed to process new password: %w", err)
	}

	// Update password
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, *req.UserID, user); err != nil {
		log.Printf("Error updating password: %v", err)
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *AuthServiceImpl) ValidateToken(ctx context.Context, tokenString string) (*dto.ValidateTokenResponse, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.ErrInvalidToken
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, apperrors.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, apperrors.ErrInvalidToken
	}

	// Check token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return nil, apperrors.ErrInvalidToken
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, apperrors.ErrInvalidToken
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, apperrors.ErrInvalidToken
	}

	username, _ := claims["username"].(string)
	role, _ := claims["role"].(string)

	// Get email from claims if present (optional)
	var email *string
	if emailStr, ok := claims["email"].(string); ok && emailStr != "" {
		email = &emailStr
	}

	return &dto.ValidateTokenResponse{
		Valid:    true,
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
	}, nil
}

func (s *AuthServiceImpl) Logout(ctx context.Context, userID uuid.UUID) error {
	// TODO: Add token blacklisting logic
	// For now, just log the logout
	log.Printf("User logged out: %s", userID)
	return nil
}

func (s *AuthServiceImpl) RevokeAllTokens(ctx context.Context, userID uuid.UUID) error {
	// TODO: Add logic to invalidate all tokens for user
	// For now, just log the action
	log.Printf("All tokens revoked for user: %s", userID)
	return nil
}

// Helper methods
func (s *AuthServiceImpl) generateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"type":     "access",
		"user_id":  user.ID.String(),
		"username": user.Username,
		"role":     string(user.Role),
		"status":   string(user.Status),
		"exp":      time.Now().Add(s.tokenExpiration).Unix(),
		"iat":      time.Now().Unix(),
		"nbf":      time.Now().Unix(),
	}

	// Add email to claims if present
	if user.Email != nil {
		claims["email"] = *user.Email
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthServiceImpl) generateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"type":    "refresh",
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}