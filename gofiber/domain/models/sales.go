package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sales struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    *uuid.UUID `gorm:"type:uuid;index:idx_sales_user_id" json:"user_id"`
	Code      string     `gorm:"uniqueIndex;not null;size:20" json:"code"`
	FirstName string     `gorm:"not null;size:100;column:first_name" json:"first_name"`
	LastName  string     `gorm:"not null;size:100;column:last_name" json:"last_name"`
	NickName  string     `gorm:"not null;size:50;column:nick_name" json:"nick_name"`
	Active    bool       `gorm:"not null;default:true" json:"active"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User     *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	Bookings []Booking `gorm:"foreignKey:SalesID" json:"bookings,omitempty"`
}

func (Sales) TableName() string {
	return "sales"
}

func (s *Sales) GetFullName() string {
	return s.FirstName + " " + s.LastName
}

func (s *Sales) GetDisplayName() string {
	return s.NickName
}

func (s *Sales) IsActive() bool {
	return s.Active
}

func (s *Sales) Deactivate() {
	s.Active = false
}

func (s *Sales) Activate() {
	s.Active = true
}

func (s *Sales) HasUser() bool {
	return s.UserID != nil
}

func (s *Sales) IsLinkedToUser(userID uuid.UUID) bool {
	return s.UserID != nil && *s.UserID == userID
}

func (s *Sales) GetUserRole() string {
	if s.User != nil {
		return string(s.User.Role)
	}
	return ""
}

func (s *Sales) CanManageBookings() bool {
	return s.Active && s.HasUser()
}