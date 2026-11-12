package main

import (
	"context"
	"fmt"
)

// VerifyPayments verifies payment data
func VerifyPayments(ctx context.Context) error {
	fmt.Println("\n=== Verifying Payment Data ===\n")

	// Count bookings with sum_payment > 0
	var bookingsWithPayment int
	err := GetPostgresDB().QueryRowContext(ctx,
		"SELECT COUNT(*) FROM bookings WHERE sum_payment > 0 AND deleted_at IS NULL").Scan(&bookingsWithPayment)
	if err != nil {
		return err
	}
	fmt.Printf("Bookings with sum_payment > 0: %d\n", bookingsWithPayment)

	// Total sum_payment
	var totalPayment float64
	err = GetPostgresDB().QueryRowContext(ctx,
		"SELECT COALESCE(SUM(sum_payment), 0) FROM bookings WHERE deleted_at IS NULL").Scan(&totalPayment)
	if err != nil {
		return err
	}
	fmt.Printf("Total sum_payment: %.2f\n", totalPayment)

	// Count payment slips
	var totalSlips int
	err = GetPostgresDB().QueryRowContext(ctx,
		"SELECT COUNT(*) FROM payment_slips WHERE deleted_at IS NULL").Scan(&totalSlips)
	if err != nil {
		return err
	}
	fmt.Printf("Total payment slips: %d\n", totalSlips)

	// Count confirmed slips
	var confirmedSlips int
	err = GetPostgresDB().QueryRowContext(ctx,
		"SELECT COUNT(*) FROM payment_slips WHERE status = 'confirm' AND deleted_at IS NULL").Scan(&confirmedSlips)
	if err != nil {
		return err
	}
	fmt.Printf("Confirmed payment slips: %d\n", confirmedSlips)

	// Sum of confirmed slip amounts
	var slipTotal float64
	err = GetPostgresDB().QueryRowContext(ctx,
		"SELECT COALESCE(SUM(amount), 0) FROM payment_slips WHERE status = 'confirm' AND deleted_at IS NULL").Scan(&slipTotal)
	if err != nil {
		return err
	}
	fmt.Printf("Total confirmed slip amount: %.2f\n", slipTotal)

	// Sample bookings with payments
	fmt.Println("\n--- Sample Bookings with Payments ---")
	rows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT b.invoice_id, b.sum_payment, b.sum_net,
			   (SELECT COUNT(*) FROM payment_slips ps WHERE ps.booking_id = b.id AND ps.status = 'confirm' AND ps.deleted_at IS NULL) as slip_count,
			   (SELECT COALESCE(SUM(amount), 0) FROM payment_slips ps WHERE ps.booking_id = b.id AND ps.status = 'confirm' AND ps.deleted_at IS NULL) as slip_total
		FROM bookings b
		WHERE b.sum_payment > 0 AND b.deleted_at IS NULL
		ORDER BY b.sum_payment DESC
		LIMIT 10
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Printf("\n%-20s %15s %15s %10s %15s\n", "Invoice ID", "sum_payment", "sum_net", "Slips", "Slip Total")
	fmt.Println("------------------------------------------------------------------------------")

	for rows.Next() {
		var invoiceID string
		var sumPayment, sumNet, slipTotal float64
		var slipCount int
		if err := rows.Scan(&invoiceID, &sumPayment, &sumNet, &slipCount, &slipTotal); err != nil {
			continue
		}
		fmt.Printf("%-20s %15.2f %15.2f %10d %15.2f\n", invoiceID, sumPayment, sumNet, slipCount, slipTotal)
	}

	return nil
}
