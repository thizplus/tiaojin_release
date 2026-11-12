package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// FixSlipAmounts fixes payment slip amounts by calculating from booking.sumPayment
func FixSlipAmounts(ctx context.Context) error {
	fmt.Println("\n=== Fixing Payment Slip Amounts ===\n")

	// Get all bookings from MongoDB with slips and sumPayment > 0
	collection := GetMongoDB().Collection("tbl_booking")

	cursor, err := collection.Find(ctx, bson.M{
		"slip":       bson.M{"$exists": true, "$ne": nil, "$not": bson.M{"$size": 0}},
		"sumPayment": bson.M{"$gt": 0},
	})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	updated := 0
	skipped := 0
	errors := 0

	for cursor.Next(ctx) {
		var booking bson.M
		if err := cursor.Decode(&booking); err != nil {
			errors++
			continue
		}

		invoiceID, _ := booking["invoiceId"].(string)
		sumPayment := toFloat64(booking["sumPayment"])

		if sumPayment <= 0 {
			skipped++
			continue
		}

		// Count confirmed slips
		slipArr, ok := booking["slip"].(bson.A)
		if !ok {
			skipped++
			continue
		}

		confirmedSlipIDs := []string{}
		for _, s := range slipArr {
			if slip, ok := s.(bson.M); ok {
				status, _ := slip["status"].(string)
				if status == "confirm" || status == "confirmed" || status == "approved" {
					if id, ok := slip["id"].(string); ok {
						confirmedSlipIDs = append(confirmedSlipIDs, id)
					}
				}
			}
		}

		if len(confirmedSlipIDs) == 0 {
			skipped++
			continue
		}

		// Calculate amount per slip
		amountPerSlip := sumPayment / float64(len(confirmedSlipIDs))

		// Get booking ID from PostgreSQL
		var bookingUUID string
		err := GetPostgresDB().QueryRowContext(ctx,
			"SELECT id FROM bookings WHERE invoice_id = $1 AND deleted_at IS NULL LIMIT 1",
			invoiceID).Scan(&bookingUUID)
		if err != nil {
			skipped++
			continue
		}

		// Update all confirmed slips for this booking
		updateQuery := `
			UPDATE payment_slips
			SET amount = $1, updated_at = NOW()
			WHERE booking_id = $2 AND status = 'confirm' AND deleted_at IS NULL
		`
		result, err := GetPostgresDB().ExecContext(ctx, updateQuery, amountPerSlip, bookingUUID)
		if err != nil {
			fmt.Printf("    ERROR updating slips for %s: %v\n", invoiceID, err)
			errors++
			continue
		}

		affected, _ := result.RowsAffected()
		if affected > 0 {
			updated += int(affected)
		}
	}

	fmt.Printf("\n=== Fix Slip Amounts Summary ===\n")
	fmt.Printf("    Updated: %d slips\n", updated)
	fmt.Printf("    Skipped: %d bookings\n", skipped)
	fmt.Printf("    Errors: %d\n", errors)

	// Also update sum_payment in bookings table based on confirmed slips
	fmt.Println("\n--- Updating booking sum_payment from confirmed slips ---")
	updateSumQuery := `
		UPDATE bookings b
		SET sum_payment = COALESCE((
			SELECT SUM(ps.amount)
			FROM payment_slips ps
			WHERE ps.booking_id = b.id
			AND ps.status = 'confirm'
			AND ps.deleted_at IS NULL
		), 0),
		updated_at = NOW()
		WHERE b.deleted_at IS NULL
	`
	result, err := GetPostgresDB().ExecContext(ctx, updateSumQuery)
	if err != nil {
		fmt.Printf("    ERROR updating sum_payment: %v\n", err)
	} else {
		affected, _ := result.RowsAffected()
		fmt.Printf("    Updated sum_payment for %d bookings\n", affected)
	}

	return nil
}
