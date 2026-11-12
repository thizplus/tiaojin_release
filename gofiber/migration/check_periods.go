package main

import (
	"context"
	"fmt"
)

// CheckPeriodDates checks if period dates are stored correctly
func CheckPeriodDates(ctx context.Context) error {
	fmt.Println("\n=== Checking Period Dates in PostgreSQL ===\n")

	// Check periods with distinct dates
	rows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT p.id, t.code, t.title, p.start_date, p.end_date, p.seat, p.status
		FROM periods p
		LEFT JOIN tours t ON p.tour_id = t.id
		WHERE p.deleted_at IS NULL
		ORDER BY p.start_date DESC
		LIMIT 20
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Printf("%-36s %-15s %-40s %-12s %-12s %6s %s\n",
		"Period ID", "Tour Code", "Tour Title", "Start Date", "End Date", "Seat", "Status")
	fmt.Println("------------------------------------------------------------------------------------------------------------------")

	for rows.Next() {
		var id, tourCode, tourTitle, status string
		var startDate, endDate interface{}
		var seat int
		if err := rows.Scan(&id, &tourCode, &tourTitle, &startDate, &endDate, &seat, &status); err != nil {
			fmt.Printf("Error scanning: %v\n", err)
			continue
		}
		// Truncate title if too long
		if len(tourTitle) > 38 {
			tourTitle = tourTitle[:38] + ".."
		}
		fmt.Printf("%-36s %-15s %-40s %-12v %-12v %6d %s\n",
			id[:8]+"...", tourCode, tourTitle, startDate, endDate, seat, status)
	}

	// Check distinct dates count
	var distinctDates int
	err = GetPostgresDB().QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT start_date) FROM periods WHERE deleted_at IS NULL
	`).Scan(&distinctDates)
	if err != nil {
		return err
	}
	fmt.Printf("\nDistinct start_date values: %d\n", distinctDates)

	// Check total periods
	var totalPeriods int
	err = GetPostgresDB().QueryRowContext(ctx, `
		SELECT COUNT(*) FROM periods WHERE deleted_at IS NULL
	`).Scan(&totalPeriods)
	if err != nil {
		return err
	}
	fmt.Printf("Total periods: %d\n", totalPeriods)

	// Check MongoDB periods for comparison
	fmt.Println("\n=== Checking Period Dates in MongoDB ===\n")

	collection := GetMongoDB().Collection("tbl_period")
	cursor, err := collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) && count < 10 {
		var period map[string]interface{}
		if err := cursor.Decode(&period); err != nil {
			continue
		}
		fmt.Printf("MongoDB Period: start_date=%v, end_date=%v\n",
			period["start_date"], period["end_date"])
		count++
	}

	return nil
}

// CleanDuplicatePeriods removes duplicate periods keeping only the first occurrence
// Duplicates are identified by tour_id + start_date + end_date
func CleanDuplicatePeriods(ctx context.Context) error {
	fmt.Println("\n=== Cleaning Duplicate Periods ===\n")

	// First, count duplicates
	var duplicateCount int
	err := GetPostgresDB().QueryRowContext(ctx, `
		SELECT COUNT(*) FROM (
			SELECT tour_id, start_date, end_date, COUNT(*) as cnt
			FROM periods
			WHERE deleted_at IS NULL
			GROUP BY tour_id, start_date, end_date
			HAVING COUNT(*) > 1
		) AS duplicates
	`).Scan(&duplicateCount)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d duplicate groups (tour_id + start_date + end_date combinations with multiple periods)\n", duplicateCount)

	if duplicateCount == 0 {
		fmt.Println("No duplicates to clean!")
		return nil
	}

	// Get total excess periods
	var excessPeriods int
	err = GetPostgresDB().QueryRowContext(ctx, `
		SELECT COALESCE(SUM(cnt - 1), 0) FROM (
			SELECT tour_id, start_date, end_date, COUNT(*) as cnt
			FROM periods
			WHERE deleted_at IS NULL
			GROUP BY tour_id, start_date, end_date
			HAVING COUNT(*) > 1
		) AS duplicates
	`).Scan(&excessPeriods)
	if err != nil {
		return err
	}
	fmt.Printf("Total excess periods to delete: %d\n", excessPeriods)

	// Delete duplicates keeping the first (by created_at)
	result, err := GetPostgresDB().ExecContext(ctx, `
		DELETE FROM periods
		WHERE id IN (
			SELECT id FROM (
				SELECT id,
					ROW_NUMBER() OVER (
						PARTITION BY tour_id, start_date, end_date
						ORDER BY created_at ASC
					) AS rn
				FROM periods
				WHERE deleted_at IS NULL
			) AS ranked
			WHERE rn > 1
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to delete duplicates: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("\nDeleted %d duplicate periods\n", rowsAffected)

	// Verify counts after cleanup
	var afterCount int
	GetPostgresDB().QueryRowContext(ctx, `SELECT COUNT(*) FROM periods WHERE deleted_at IS NULL`).Scan(&afterCount)
	fmt.Printf("Periods remaining: %d\n", afterCount)

	return nil
}
