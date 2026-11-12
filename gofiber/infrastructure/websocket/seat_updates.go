package websocket

import (
	"context"
	"fmt"
	"log"
	"time"

	"gofiber-template/domain/services"
	"github.com/google/uuid"
)

type SeatUpdateData struct {
	PeriodID       string `json:"period_id"`
	TotalSeats     int    `json:"total_seats"`
	BookedSeats    int    `json:"booked_seats"`
	HeldSeats      int    `json:"held_seats"`
	AvailableSeats int    `json:"available_seats"`
	UpdateType     string `json:"update_type"` // "booking_created", "booking_updated", "seat_held", "seat_released"
	Timestamp      string `json:"timestamp"`
}

func BroadcastSeatUpdate(ctx context.Context, periodID uuid.UUID, updateType string, seatService services.SeatAvailabilityService) {
	// Calculate current seat availability
	availability, err := seatService.CalculateAvailability(periodID)
	if err != nil {
		log.Printf("Failed to calculate seat availability for broadcast: %v", err)
		return
	}

	// Create update data
	updateData := SeatUpdateData{
		PeriodID:       periodID.String(),
		TotalSeats:     availability.TotalSeats,
		BookedSeats:    availability.BookedSeats,
		HeldSeats:      availability.HeldSeats,
		AvailableSeats: availability.AvailableSeats,
		UpdateType:     updateType,
		Timestamp:      getCurrentTimestamp(),
	}

	// Broadcast to all clients (global seat updates)
	Manager.BroadcastToAll("seat_update", updateData)

	// Broadcast to period-specific room
	roomID := fmt.Sprintf("period_%s", periodID.String())
	Manager.BroadcastToRoom(roomID, "seat_update", updateData)

	log.Printf("Broadcasted seat update for period %s: available=%d, type=%s",
		periodID.String(), availability.AvailableSeats, updateType)
}

func getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}