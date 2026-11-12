package serviceimpl

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
)

type SeatAvailabilityServiceImpl struct {
	periodRepo   repositories.PeriodRepository
	bookingRepo  repositories.BookingRepository
	seatHoldRepo repositories.SeatHoldRepository
}

func NewSeatAvailabilityService(
	periodRepo repositories.PeriodRepository,
	bookingRepo repositories.BookingRepository,
	seatHoldRepo repositories.SeatHoldRepository,
) services.SeatAvailabilityService {
	return &SeatAvailabilityServiceImpl{
		periodRepo:   periodRepo,
		bookingRepo:  bookingRepo,
		seatHoldRepo: seatHoldRepo,
	}
}

// CalculateAvailability คำนวณจำนวนที่นั่งว่างสำหรับ period
// ใช้วิธีเดียวกับ admin periods route - ใช้ field Booking จาก Period model และนับ seat holds
func (s *SeatAvailabilityServiceImpl) CalculateAvailability(periodID uuid.UUID) (*services.SeatAvailability, error) {
	log.Printf("🔍 Calculating seat availability for period: %s", periodID)

	// 1. ดึงข้อมูล Period
	ctx := context.Background()
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		return nil, fmt.Errorf("failed to get period: %w", err)
	}

	// 2. นับจำนวน bookings ที่ active (เหมือน admin periods route - ใช้ real-time calculation)
	bookedSeats, err := s.countConfirmedBookings(periodID)
	if err != nil {
		log.Printf("⚠️ Warning: Failed to count confirmed bookings: %v", err)
		bookedSeats = 0 // ถ้าไม่สามารถนับได้ให้ใช้ 0
	}

	// 3. นับจำนวน seat holds ที่ยัง active
	activeSeatHolds, err := s.countActiveSeatHolds(periodID)
	if err != nil {
		log.Printf("⚠️ Warning: Failed to count active seat holds: %v", err)
		activeSeatHolds = 0 // ถ้าไม่สามารถนับได้ให้ใช้ 0
	}

	// 4. คำนวณที่นั่งว่าง (ตาม logic ของ admin periods route)
	availableSeats := period.Seat - bookedSeats - activeSeatHolds
	if availableSeats < 0 {
		availableSeats = 0 // ป้องกันไม่ให้เป็นค่าลบ
	}

	result := &services.SeatAvailability{
		PeriodID:       periodID,
		TotalSeats:     period.Seat,
		BookedSeats:    bookedSeats,
		HeldSeats:      activeSeatHolds,
		AvailableSeats: availableSeats,
	}

	log.Printf("✅ Seat availability calculated (using real-time booking count): Total=%d, Booked=%d, Held=%d, Available=%d",
		result.TotalSeats, result.BookedSeats, result.HeldSeats, result.AvailableSeats)

	return result, nil
}

// countConfirmedBookings นับจำนวน bookings ที่ active (เหมือน admin periods route)
func (s *SeatAvailabilityServiceImpl) countConfirmedBookings(periodID uuid.UUID) (int, error) {
	// ใช้สถานะเดียวกับ admin periods route: pending, processing, confirm
	activeStatuses := []string{
		string(models.BookingStatusPending),
		string(models.BookingStatusProcessing),
		string(models.BookingStatusConfirm),
	}

	bookings, err := s.bookingRepo.FindByPeriodIDAndStatuses(periodID, activeStatuses)
	if err != nil {
		return 0, err
	}

	totalSeats := 0
	for _, booking := range bookings {
		totalSeats += booking.SumQty
	}

	log.Printf("📊 Active bookings for period %s: %d bookings, %d total seats (pending+processing+confirm)",
		periodID, len(bookings), totalSeats)

	return totalSeats, nil
}

// countActiveSeatHolds นับจำนวน seat holds ที่ยัง active
func (s *SeatAvailabilityServiceImpl) countActiveSeatHolds(periodID uuid.UUID) (int, error) {
	// ดึง seat holds ที่ยัง active และยังไม่หมดอายุ
	seatHolds, err := s.seatHoldRepo.FindActiveBYPeriodID(periodID)
	if err != nil {
		return 0, err
	}

	totalHeldSeats := 0
	now := time.Now()

	for _, hold := range seatHolds {
		// ตรวจสอบว่า hold ยังไม่หมดอายุ
		if hold.ExpiresAt == nil || hold.ExpiresAt.After(now) {
			totalHeldSeats += hold.HeldSeats
		}
	}

	log.Printf("📊 Active seat holds for period %s: %d holds, %d total held seats",
		periodID, len(seatHolds), totalHeldSeats)

	return totalHeldSeats, nil
}

// GetAvailabilityForPeriods ดึงข้อมูลที่นั่งว่างสำหรับหลาย periods
func (s *SeatAvailabilityServiceImpl) GetAvailabilityForPeriods(periodIDs []uuid.UUID) (map[uuid.UUID]*services.SeatAvailability, error) {
	result := make(map[uuid.UUID]*services.SeatAvailability)

	for _, periodID := range periodIDs {
		availability, err := s.CalculateAvailability(periodID)
		if err != nil {
			log.Printf("❌ Failed to calculate availability for period %s: %v", periodID, err)
			continue
		}
		result[periodID] = availability
	}

	return result, nil
}