package main

import (
	"context"
	"fmt"
	"time"
)

// FixAll runs all fix scripts in the correct order
// This repairs all mapping issues after separate migrations
func FixAll(ctx context.Context) error {
	fmt.Println("\n============================================")
	fmt.Println("  Fix All - Repair All Mapping Issues")
	fmt.Println("============================================")
	fmt.Println()

	LogSection("FIX ALL - START")
	startTime := time.Now()

	// Fix scripts in order of dependencies
	// Order matters: Cities → Tours → Periods → Users → Sales → Bookings
	// NOTE: FixPeriodTours is REMOVED because sync-full already sets correct tour_id
	// Running it would incorrectly reassign periods based on date matching only
	fixes := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"Tour-City Relationships", FixTourCities},          // 1. Tours need correct city_id
		// REMOVED: {"Period-Tour Relationships", FixPeriodTours} - handled by sync
		{"User-Company Relationships", FixUserCompany},      // 2. Users need company profiles
		{"User Status", FixUserStatus},                      // 3. Users need correct status
		{"Booking-User Relationships", FixBookingUsers},     // 4. Bookings need correct user_id
		{"Booking-Tour Relationships", FixBookingTours},     // 5. Bookings need correct tour_id
		{"Booking-Period Relationships", FixBookingPeriods}, // 6. Bookings need correct period_id
		{"Booking-Sales Relationships", FixBookingSales},    // 7. Bookings need correct sales_id
		{"Payment Slip Amounts", FixSlipAmounts},            // 8. Fix slip amounts from sumPayment
	}

	for i, fix := range fixes {
		fmt.Printf("=== [%d/%d] %s ===\n", i+1, len(fixes), fix.name)
		LogFixStart(fix.name)
		fixStart := time.Now()

		if err := fix.fn(ctx); err != nil {
			LogError("Failed at %s: %v", fix.name, err)
			return fmt.Errorf("failed at %s: %w", fix.name, err)
		}

		duration := time.Since(fixStart).Round(time.Millisecond)
		fmt.Printf("    Completed in %v\n\n", duration)
		LogSuccess("%s completed in %v", fix.name, duration)
	}

	totalDuration := time.Since(startTime).Round(time.Second)
	fmt.Println("============================================")
	fmt.Println("  Fix All Complete!")
	fmt.Println("============================================")
	fmt.Printf("  Total time: %v\n", totalDuration)
	fmt.Println("============================================")

	LogSection("FIX ALL - COMPLETE")
	LogSuccess("All fixes completed in %v", totalDuration)

	return nil
}

// SyncLatest syncs latest data from MongoDB and fixes all mappings
// This is the recommended command for daily sync
func SyncLatest(ctx context.Context, since time.Time) error {
	fmt.Println("\n============================================")
	fmt.Println("  Sync Latest - Full Sync with Fixes")
	fmt.Println("============================================")
	fmt.Printf("  Syncing data since: %s\n", since.Format("2006-01-02 15:04:05"))
	fmt.Println()

	startTime := time.Now()

	// Step 1: Run incremental sync
	fmt.Println("=== Step 1: Incremental Data Sync ===")
	if err := SyncDaily(ctx, since); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}
	fmt.Println()

	// Step 2: Fix all mappings
	fmt.Println("=== Step 2: Fix All Mappings ===")
	if err := FixAll(ctx); err != nil {
		return fmt.Errorf("fix failed: %w", err)
	}

	fmt.Println("\n============================================")
	fmt.Println("  Sync Latest Complete!")
	fmt.Println("============================================")
	fmt.Printf("  Total time: %v\n", time.Since(startTime).Round(time.Second))
	fmt.Println("============================================")

	return nil
}
