package serviceimpl

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type UserServiceImpl struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewUserService(userRepo repositories.UserRepository, jwtSecret string) services.UserService {
	return &UserServiceImpl{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if user with username already exists
	existingUser, _ := s.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("user with this username already exists")
	}

	// Check if email already exists (if provided)
	if req.Email != nil && *req.Email != "" {
		existingUser, _ := s.userRepo.GetByEmail(ctx, *req.Email)
		if existingUser != nil {
			return nil, errors.New("user with this email already exists")
		}
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = "agent" // Default role
	}

	// Create user model
	user := &models.User{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      models.UserRole(role),
		Status:    models.UserStatusPending, // Default status for admin-created users
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Return user response (without password)
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(user.Role),
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return dto.UserToUserResponse(user), nil
}

func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return dto.UserToUserResponse(user), nil
}

func (s *UserServiceImpl) GetUserWithProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	// First get user to find username
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if user is nil
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Then get user with profile loaded
	userWithProfile, err := s.userRepo.GetByUsernameWithProfile(ctx, user.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return dto.UserToUserResponse(userWithProfile), nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, userID uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Apply updates
	if req.Email != nil {
		user.Email = req.Email
	}
	if req.Role != nil {
		user.Role = models.UserRole(*req.Role)
	}
	if req.Status != nil {
		user.Status = models.UserStatus(*req.Status)
	}

	err = s.userRepo.Update(ctx, userID, user)
	if err != nil {
		return nil, err
	}

	return dto.UserToUserResponse(user), nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.Delete(ctx, userID)
}

func (s *UserServiceImpl) UpdateUserStatus(ctx context.Context, userID uuid.UUID, req dto.UpdateUserStatusRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = s.userRepo.UpdateStatus(ctx, userID, models.UserStatus(req.Status))
	if err != nil {
		return nil, err
	}

	// Get updated user
	user, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.UserToUserResponse(user), nil
}

func (s *UserServiceImpl) ChangeUserRole(ctx context.Context, userID uuid.UUID, req dto.ChangeUserRoleRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update user role
	user.Role = models.UserRole(req.Role)
	user.UpdatedAt = time.Now()

	err = s.userRepo.Update(ctx, userID, user)
	if err != nil {
		return nil, err
	}

	// Get updated user
	user, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.UserToUserResponse(user), nil
}

func (s *UserServiceImpl) ListUsers(ctx context.Context, req dto.UserListRequest) (*dto.UserListResponse, error) {
	var users []*models.User
	var total int64
	var err error

	// Handle search if provided (search by username)
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		users, err = s.userRepo.SearchByUsername(ctx, strings.TrimSpace(*req.Search), req.Offset, req.Limit)
		if err != nil {
			return nil, err
		}
		// For search, approximate total
		total = int64(len(users))
	} else {
		users, err = s.userRepo.List(ctx, req.Offset, req.Limit)
		if err != nil {
			return nil, err
		}

		total, err = s.userRepo.Count(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Convert to response DTOs
	var userResponses []*dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserToUserResponse(user))
	}

	pagination := utils.CalculatePagination(total, req.Offset, req.Limit)

	return &dto.UserListResponse{
		Users:      userResponses,
		Pagination: pagination,
	}, nil
}

func (s *UserServiceImpl) GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"role":     string(user.Role),
		"status":   string(user.Status),
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":      time.Now().Unix(),
	}

	// Add email to claims if present
	if user.Email != nil {
		claims["email"] = *user.Email
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserServiceImpl) ValidateJWT(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("invalid token claims")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, errors.New("invalid user ID in token")
		}

		user, err := s.userRepo.GetByID(context.Background(), userID)
		if err != nil {
			return nil, errors.New("user not found")
		}

		return user, nil
	}

	return nil, errors.New("invalid token")
}

func (s *UserServiceImpl) GetUserStats(ctx context.Context) (*dto.UserStatsResponse, error) {
	totalUsers, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// Count users by status
	activeCount, err := s.userRepo.CountByStatus(ctx, models.UserStatusActive)
	if err != nil {
		return nil, err
	}

	pendingCount, err := s.userRepo.CountByStatus(ctx, models.UserStatusPending)
	if err != nil {
		return nil, err
	}

	inactiveCount, err := s.userRepo.CountByStatus(ctx, models.UserStatusInactive)
	if err != nil {
		return nil, err
	}

	return &dto.UserStatsResponse{
		TotalUsers:    totalUsers,
		ActiveUsers:   activeCount,
		PendingUsers:  pendingCount,
		InactiveUsers: inactiveCount,
	}, nil
}

func (s *UserServiceImpl) ApproveUser(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	err := s.userRepo.UpdateStatus(ctx, userID, models.UserStatusActive)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.UserToUserResponse(user), nil
}

func (s *UserServiceImpl) SuspendUser(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	err := s.userRepo.UpdateStatus(ctx, userID, models.UserStatusInactive)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.UserToUserResponse(user), nil
}

func (s *UserServiceImpl) ReactivateUser(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	err := s.userRepo.UpdateStatus(ctx, userID, models.UserStatusActive)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.UserToUserResponse(user), nil
}
