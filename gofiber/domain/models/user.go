package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleSuper UserRole = "super" // เจ้าของระบบ - full access
	UserRoleAdmin UserRole = "admin" // ผู้ดูแล - management functions
	UserRoleAgent UserRole = "agent" // ตัวแทนขาย - booking management, customer service
	UserRoleSales UserRole = "sales" // คนขายให้ agent - commission tracking, sales operations
)

type UserStatus string

const (
	UserStatusPending  UserStatus = "pending"
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

type User struct {
	ID       uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Username string     `gorm:"uniqueIndex;not null;size:100" json:"username"`
	Email    *string    `gorm:"uniqueIndex;size:255" json:"email,omitempty"` // Optional, สำหรับ reset password
	Password string     `gorm:"not null;size:255" json:"-"`
	Role     UserRole   `gorm:"type:varchar(20);not null;default:'agent'" json:"role"`
	Status   UserStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	CompanyProfile *CompanyProfile `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"company_profile,omitempty"`
	Bookings       []Booking       `gorm:"foreignKey:UserID" json:"bookings,omitempty"`
}

func (User) TableName() string {
	return "users"
}