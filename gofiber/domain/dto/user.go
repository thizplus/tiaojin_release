package dto

import "github.com/google/uuid"

type CreateUserRequest struct {
	Username string  `json:"username" validate:"required,min=3,max=100"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email,max=255"`
	Password string  `json:"password" validate:"required,min=8,max=72"`
	Role     string  `json:"role" validate:"omitempty,oneof=super admin agent sales"`
}

type UserResponse struct {
	ID             uuid.UUID               `json:"id"`
	Username       string                  `json:"username"`
	Email          *string                 `json:"email,omitempty"`
	Role           string                  `json:"role"`
	Status         string                  `json:"status"`
	CreatedAt      string                  `json:"created_at"`
	UpdatedAt      string                  `json:"updated_at"`
	CompanyProfile *CompanyProfileResponse `json:"company_profile,omitempty"`
}

type UpdateUserRequest struct {
	Email  *string `json:"email,omitempty" validate:"omitempty,email"`
	Role   *string `json:"role,omitempty" validate:"omitempty,oneof=super admin agent sales"`
	Status *string `json:"status,omitempty" validate:"omitempty,oneof=pending active inactive"`
}

type UpdateUserStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending active inactive"`
}

type UserListRequest struct {
	Status *string `json:"status,omitempty" validate:"omitempty,oneof=pending active inactive"`
	Role   *string `json:"role,omitempty" validate:"omitempty,oneof=super admin agent sales"`
	Search *string `json:"search,omitempty"`
	Offset int     `json:"offset" validate:"min=0"`
	Limit  int     `json:"limit" validate:"min=1,max=100"`
}

type UserListResponse struct {
	Users      []*UserResponse `json:"users"`
	Pagination PaginationMeta  `json:"pagination"`
}

type UserStatsResponse struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	PendingUsers  int64 `json:"pending_users"`
	InactiveUsers int64 `json:"inactive_users"`
}

type UserProfileRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserProfileRequest struct {
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,max=50"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,max=50"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,max=20"`
}

type ChangeUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=super admin agent sales"`
}