package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeatHoldStatus represents the status of a seat hold
type SeatHoldStatus string

const (
	SeatHoldStatusActive    SeatHoldStatus = "active"    // Hold is currently active
	SeatHoldStatusReleased  SeatHoldStatus = "released"  // Hold was manually released
	SeatHoldStatusExpired   SeatHoldStatus = "expired"   // Hold expired automatically
	SeatHoldStatusConverted SeatHoldStatus = "converted" // Hold was converted to booking
)

// SeatHold represents a temporary seat reservation by admin
type SeatHold struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PeriodID  uuid.UUID      `json:"period_id" gorm:"type:uuid;not null;index"`
	AdminID   uuid.UUID      `json:"admin_id" gorm:"type:uuid;not null;index"`
	HeldSeats int            `json:"held_seats" gorm:"not null;check:held_seats > 0"`
	Reason    *string        `json:"reason" gorm:"type:text"`
	ExpiresAt *time.Time     `json:"expires_at" gorm:"index"`
	Status    SeatHoldStatus `json:"status" gorm:"type:varchar(20);not null;default:'active';index"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null;default:now()"`

	// Relations
	Period *Period `json:"period,omitempty" gorm:"foreignKey:PeriodID;constraint:OnDelete:CASCADE"`
	Admin  *User   `json:"admin,omitempty" gorm:"foreignKey:AdminID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for SeatHold
func (SeatHold) TableName() string {
	return "seat_holds"
}

// BeforeCreate hook
func (sh *SeatHold) BeforeCreate(tx *gorm.DB) error {
	if sh.ID == uuid.Nil {
		sh.ID = uuid.New()
	}
	return nil
}

// IsActive returns true if the hold is currently active
func (sh *SeatHold) IsActive() bool {
	return sh.Status == SeatHoldStatusActive
}

// IsExpired returns true if the hold has an expiry time and it has passed
func (sh *SeatHold) IsExpired() bool {
	if sh.ExpiresAt == nil {
		return false // No expiry time set
	}
	return time.Now().After(*sh.ExpiresAt)
}

// ShouldAutoExpire returns true if this hold should be auto-expired
func (sh *SeatHold) ShouldAutoExpire() bool {
	return sh.IsActive() && sh.IsExpired()
}

// GetTimeRemaining returns the time remaining until expiry, or nil if no expiry set
func (sh *SeatHold) GetTimeRemaining() *time.Duration {
	if sh.ExpiresAt == nil {
		return nil
	}
	remaining := time.Until(*sh.ExpiresAt)
	return &remaining
}

// Expire marks the hold as expired
func (sh *SeatHold) Expire() {
	sh.Status = SeatHoldStatusExpired
	sh.UpdatedAt = time.Now()
}

// Release marks the hold as released
func (sh *SeatHold) Release() {
	sh.Status = SeatHoldStatusReleased
	sh.UpdatedAt = time.Now()
}

// Convert marks the hold as converted to booking
func (sh *SeatHold) Convert() {
	sh.Status = SeatHoldStatusConverted
	sh.UpdatedAt = time.Now()
}

// ExtendExpiry extends the hold expiry time
func (sh *SeatHold) ExtendExpiry(duration time.Duration) {
	if sh.ExpiresAt == nil {
		newExpiry := time.Now().Add(duration)
		sh.ExpiresAt = &newExpiry
	} else {
		newExpiry := sh.ExpiresAt.Add(duration)
		sh.ExpiresAt = &newExpiry
	}
	sh.UpdatedAt = time.Now()
}

// SetExpiry sets a new expiry time
func (sh *SeatHold) SetExpiry(expiryTime time.Time) {
	sh.ExpiresAt = &expiryTime
	sh.UpdatedAt = time.Now()
}

// RemoveExpiry removes the expiry time (makes hold permanent until manually released)
func (sh *SeatHold) RemoveExpiry() {
	sh.ExpiresAt = nil
	sh.UpdatedAt = time.Now()
}

// ValidateStatus validates if the status transition is allowed
func (sh *SeatHold) ValidateStatus(newStatus SeatHoldStatus) bool {
	currentStatus := sh.Status

	// Define allowed transitions
	allowedTransitions := map[SeatHoldStatus][]SeatHoldStatus{
		SeatHoldStatusActive: {
			SeatHoldStatusReleased,
			SeatHoldStatusExpired,
			SeatHoldStatusConverted,
		},
		SeatHoldStatusReleased:  {}, // Final state
		SeatHoldStatusExpired:   {}, // Final state
		SeatHoldStatusConverted: {}, // Final state
	}

	allowed, exists := allowedTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, allowedStatus := range allowed {
		if allowedStatus == newStatus {
			return true
		}
	}

	return false
}

// SeatHoldSummary provides aggregated information about seat holds for a period
type SeatHoldSummary struct {
	PeriodID          uuid.UUID `json:"period_id"`
	TotalActiveHolds  int       `json:"total_active_holds"`
	TotalHeldSeats    int       `json:"total_held_seats"`
	TotalExpiredHolds int       `json:"total_expired_holds"`
	EarliestExpiry    *time.Time `json:"earliest_expiry,omitempty"`
	LatestExpiry      *time.Time `json:"latest_expiry,omitempty"`
}