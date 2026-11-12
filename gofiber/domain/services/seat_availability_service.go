package services

import (
	"github.com/google/uuid"
)

// SeatAvailability represents the current seat status for a period
type SeatAvailability struct {
	PeriodID       uuid.UUID `json:"period_id"`
	TotalSeats     int       `json:"total_seats"`
	BookedSeats    int       `json:"booked_seats"`
	HeldSeats      int       `json:"held_seats"`
	AvailableSeats int       `json:"available_seats"`
}

// SeatAvailabilityService handles seat availability calculations
type SeatAvailabilityService interface {
	// CalculateAvailability คำนวณจำนวนที่นั่งว่างสำหรับ period
	CalculateAvailability(periodID uuid.UUID) (*SeatAvailability, error)

	// GetAvailabilityForPeriods ดึงข้อมูลที่นั่งว่างสำหรับหลาย periods
	GetAvailabilityForPeriods(periodIDs []uuid.UUID) (map[uuid.UUID]*SeatAvailability, error)
}