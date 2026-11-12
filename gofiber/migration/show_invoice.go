package main

import (
	"context"
	"fmt"
)

// ShowInvoice shows detailed invoice data including passports and slips
func ShowInvoice(ctx context.Context, invoiceID string) error {
	fmt.Printf("\n=== Invoice Details: %s ===\n\n", invoiceID)

	// Get booking details
	var bookingID, customerName string
	var tourName, periodStartDate interface{}
	var sumPayment, sumNet float64
	err := GetPostgresDB().QueryRowContext(ctx, `
		SELECT b.id, COALESCE(t.title, '-'), p.start_date, b.customer_name, b.sum_payment, b.sum_net
		FROM bookings b
		LEFT JOIN tours t ON b.tour_id = t.id
		LEFT JOIN periods p ON b.period_id = p.id
		WHERE b.invoice_id = $1 AND b.deleted_at IS NULL
	`, invoiceID).Scan(&bookingID, &tourName, &periodStartDate, &customerName, &sumPayment, &sumNet)
	if err != nil {
		return fmt.Errorf("booking not found: %w", err)
	}

	fmt.Printf("Customer: %s\n", customerName)
	fmt.Printf("Tour: %v\n", tourName)
	fmt.Printf("Period Start: %v\n", periodStartDate)
	fmt.Printf("Sum Net: %.2f\n", sumNet)
	fmt.Printf("Sum Payment: %.2f\n", sumPayment)

	// Get payment slips
	fmt.Println("\n--- Payment Slips ---")
	slipRows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT id, amount, status, COALESCE(image, ''), created_at
		FROM payment_slips
		WHERE booking_id = $1 AND deleted_at IS NULL
		ORDER BY created_at
	`, bookingID)
	if err != nil {
		return err
	}
	defer slipRows.Close()

	slipCount := 0
	for slipRows.Next() {
		var id string
		var amount float64
		var status, image string
		var createdAt interface{}
		if err := slipRows.Scan(&id, &amount, &status, &image, &createdAt); err != nil {
			continue
		}
		slipCount++
		fmt.Printf("%d. Amount: %.2f, Status: %s, Image: %s\n", slipCount, amount, status, image)
	}
	if slipCount == 0 {
		fmt.Println("(No payment slips)")
	}

	// Get passports
	fmt.Println("\n--- Passenger Passports ---")
	passportRows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT id, COALESCE(image, ''), passport_number, created_at
		FROM passenger_passports
		WHERE booking_id = $1 AND deleted_at IS NULL
		ORDER BY created_at
	`, bookingID)
	if err != nil {
		return err
	}
	defer passportRows.Close()

	passCount := 0
	for passportRows.Next() {
		var id, image, passportNumber string
		var createdAt interface{}
		if err := passportRows.Scan(&id, &image, &passportNumber, &createdAt); err != nil {
			continue
		}
		passCount++
		fmt.Printf("%d. Passport: %s, Image: %s\n", passCount, passportNumber, image)
	}
	if passCount == 0 {
		fmt.Println("(No passports)")
	}

	return nil
}

// ListInvoicesWithPassports lists invoices that have passports
func ListInvoicesWithPassports(ctx context.Context) error {
	fmt.Println("\n=== Bookings with Passports ===\n")

	rows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT b.invoice_id, b.customer_name, b.sum_payment,
			   (SELECT COUNT(*) FROM payment_slips ps WHERE ps.booking_id = b.id AND ps.deleted_at IS NULL) as slip_count,
			   (SELECT COUNT(*) FROM passenger_passports pp WHERE pp.booking_id = b.id AND pp.deleted_at IS NULL) as passport_count
		FROM bookings b
		WHERE b.deleted_at IS NULL
		AND EXISTS (SELECT 1 FROM passenger_passports pp WHERE pp.booking_id = b.id AND pp.deleted_at IS NULL)
		ORDER BY b.created_at DESC
		LIMIT 10
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Printf("%-20s %-30s %15s %8s %10s\n", "Invoice ID", "Customer", "Sum Payment", "Slips", "Passports")
	fmt.Println("--------------------------------------------------------------------------------------")

	for rows.Next() {
		var invoiceID, customerName string
		var sumPayment float64
		var slipCount, passportCount int
		if err := rows.Scan(&invoiceID, &customerName, &sumPayment, &slipCount, &passportCount); err != nil {
			continue
		}
		// Truncate customer name if too long
		if len(customerName) > 28 {
			customerName = customerName[:28] + ".."
		}
		fmt.Printf("%-20s %-30s %15.2f %8d %10d\n", invoiceID, customerName, sumPayment, slipCount, passportCount)
	}

	// Count totals
	var totalBookings, totalPassports int
	GetPostgresDB().QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT b.id), COUNT(pp.id)
		FROM bookings b
		LEFT JOIN passenger_passports pp ON pp.booking_id = b.id AND pp.deleted_at IS NULL
		WHERE b.deleted_at IS NULL
	`).Scan(&totalBookings, &totalPassports)

	fmt.Printf("\nTotal: %d bookings, %d passports\n", totalBookings, totalPassports)

	return nil
}
