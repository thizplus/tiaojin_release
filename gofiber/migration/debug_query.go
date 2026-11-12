package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DebugQuery runs debug queries
func DebugQuery(ctx context.Context, query string) error {
	fmt.Println("\n=== Debug Query ===")
	fmt.Printf("Query: %s\n\n", query)

	rows, err := GetPostgresDB().QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	fmt.Printf("Columns: %v\n\n", cols)

	for rows.Next() {
		// Generic scan
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		for i, col := range cols {
			fmt.Printf("  %s: %v\n", col, values[i])
		}
		fmt.Println()
	}

	return nil
}

// CheckDuplicateTours checks for duplicate tour codes in PostgreSQL
func CheckDuplicateTours(ctx context.Context) error {
	fmt.Println("\n=== Checking Duplicate Tour Codes in PostgreSQL ===")

	rows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT code, COUNT(*) as cnt
		FROM tours
		WHERE deleted_at IS NULL
		GROUP BY code
		HAVING COUNT(*) > 1
		ORDER BY cnt DESC
		LIMIT 20
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var code string
		var cnt int
		rows.Scan(&code, &cnt)
		count++
		fmt.Printf("  %s: %d tours\n", code, cnt)
	}

	if count == 0 {
		fmt.Println("  No duplicate tour codes found")
	}

	// Check BCCKG25-SL specifically
	fmt.Println("\n--- Tours with code BCCKG25-SL ---")
	rows2, err := GetPostgresDB().QueryContext(ctx, `
		SELECT id, code, title, created_at
		FROM tours
		WHERE code = 'BCCKG25-SL' AND deleted_at IS NULL
	`)
	if err != nil {
		return err
	}
	defer rows2.Close()

	for rows2.Next() {
		var id, code, title string
		var createdAt interface{}
		rows2.Scan(&id, &code, &title, &createdAt)
		fmt.Printf("  ID: %s\n", id[:8]+"...")
		fmt.Printf("  Code: %s, Title: %s\n", code, title)
		fmt.Printf("  Created: %v\n\n", createdAt)
	}

	// Check periods for this tour
	fmt.Println("--- All periods linked to tour code BCCKG25-SL ---")
	rows3, err := GetPostgresDB().QueryContext(ctx, `
		SELECT p.id, p.tour_id, t.code, p.start_date, p.created_at
		FROM periods p
		JOIN tours t ON p.tour_id = t.id
		WHERE t.code = 'BCCKG25-SL' AND p.deleted_at IS NULL
		ORDER BY p.created_at
	`)
	if err != nil {
		return err
	}
	defer rows3.Close()

	for rows3.Next() {
		var id, tourID, code string
		var startDate, createdAt interface{}
		rows3.Scan(&id, &tourID, &code, &startDate, &createdAt)
		fmt.Printf("  Period: %s, Tour: %s, Start: %v, Created: %v\n",
			id[:8]+"...", tourID[:8]+"...", startDate, createdAt)
	}

	return nil
}

// CheckPeriodDetails checks full period data for a specific tour
func CheckPeriodDetails(ctx context.Context, tourCode string) error {
	// First, show MongoDB data
	fmt.Println("\n--- MongoDB Period Data ---")
	mongoCollection := GetMongoDB().Collection("tbl_period")

	// Find tour first
	tourCollection := GetMongoDB().Collection("tbl_tours")
	var tour struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err := tourCollection.FindOne(ctx, bson.M{"code": tourCode}).Decode(&tour)
	if err != nil {
		fmt.Printf("Tour not found in MongoDB: %v\n", err)
	} else {
		// Find period for this tour
		cursor, err := mongoCollection.Find(ctx, bson.M{"tour": tour.ID.Hex()})
		if err != nil {
			fmt.Printf("Error finding periods: %v\n", err)
		} else {
			defer cursor.Close(ctx)
			for cursor.Next(ctx) {
				var rawPeriod bson.M
				cursor.Decode(&rawPeriod)
				fmt.Printf("MongoDB Period Fields:\n")
				for key, value := range rawPeriod {
					fmt.Printf("  %s: %v\n", key, value)
				}
			}
		}
	}

	fmt.Printf("\n=== Checking Period Details for Tour: %s ===\n", tourCode)

	// Get PostgreSQL period with all fields
	fmt.Println("\n--- PostgreSQL Period Data ---")
	rows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT
			p.id, p.tour_id, p.start_date, p.end_date, p.seat,
			p.deposit, p.deposit_type, p.airline, p.airport, p.note,
			p.sector_1, p.sector_2,
			p.rate_type, p.rate_items,
			p.commission_type, p.commission_items,
			p.room_types, p.additional_services, p.files,
			p.status, p.created_at
		FROM periods p
		JOIN tours t ON p.tour_id = t.id
		WHERE t.code = $1 AND p.deleted_at IS NULL
		LIMIT 1
	`, tourCode)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, tourID string
		var startDate, endDate interface{}
		var seat int
		var deposit float64
		var depositType, airline, airport string
		var note interface{}
		var sector1, sector2 []byte
		var rateType string
		var rateItems []byte
		var commissionType string
		var commissionItems []byte
		var roomTypes, additionalServices, files []byte
		var status string
		var createdAt interface{}

		err := rows.Scan(
			&id, &tourID, &startDate, &endDate, &seat,
			&deposit, &depositType, &airline, &airport, &note,
			&sector1, &sector2,
			&rateType, &rateItems,
			&commissionType, &commissionItems,
			&roomTypes, &additionalServices, &files,
			&status, &createdAt,
		)
		if err != nil {
			fmt.Printf("Scan error: %v\n", err)
			continue
		}

		fmt.Printf("ID: %s\n", id[:8]+"...")
		fmt.Printf("Tour ID: %s\n", tourID[:8]+"...")
		fmt.Printf("Start Date: %v\n", startDate)
		fmt.Printf("End Date: %v\n", endDate)
		fmt.Printf("Seat: %d\n", seat)
		fmt.Printf("Deposit: %.2f\n", deposit)
		fmt.Printf("Deposit Type: %s\n", depositType)
		fmt.Printf("Airline: %s\n", airline)
		fmt.Printf("Airport: %s\n", airport)
		fmt.Printf("Note: %v\n", note)
		fmt.Printf("Sector 1: %s\n", string(sector1))
		fmt.Printf("Sector 2: %s\n", string(sector2))
		fmt.Printf("Rate Type: %s\n", rateType)
		fmt.Printf("Rate Items: %s\n", string(rateItems))
		fmt.Printf("Commission Type: %s\n", commissionType)
		fmt.Printf("Commission Items: %s\n", string(commissionItems))
		fmt.Printf("Room Types: %s\n", string(roomTypes))
		fmt.Printf("Additional Services: %s\n", string(additionalServices))
		fmt.Printf("Files: %s\n", string(files))
		fmt.Printf("Status: %s\n", status)
		fmt.Printf("Created At: %v\n", createdAt)

		// Check separate tables for this period
		fmt.Println("\n--- Separate Tables Data ---")

		// Check period_rate_items
		var rateItemCount int
		GetPostgresDB().QueryRowContext(ctx, `SELECT COUNT(*) FROM period_rate_items WHERE period_id = $1`, id).Scan(&rateItemCount)
		fmt.Printf("period_rate_items count: %d\n", rateItemCount)

		rateRows, _ := GetPostgresDB().QueryContext(ctx, `
			SELECT title, regular, sale, exclude, sort_order
			FROM period_rate_items
			WHERE period_id = $1
			ORDER BY sort_order`, id)
		if rateRows != nil {
			defer rateRows.Close()
			for rateRows.Next() {
				var title string
				var regular, sale float64
				var exclude bool
				var sortOrder int
				rateRows.Scan(&title, &regular, &sale, &exclude, &sortOrder)
				fmt.Printf("  - %s: regular=%.0f, sale=%.0f, exclude=%v\n", title, regular, sale, exclude)
			}
		}

		// Check period_room_types
		var roomTypeCount int
		GetPostgresDB().QueryRowContext(ctx, `SELECT COUNT(*) FROM period_room_types WHERE period_id = $1`, id).Scan(&roomTypeCount)
		fmt.Printf("\nperiod_room_types count: %d\n", roomTypeCount)

		roomRows, _ := GetPostgresDB().QueryContext(ctx, `
			SELECT title, price, sort_order
			FROM period_room_types
			WHERE period_id = $1
			ORDER BY sort_order`, id)
		if roomRows != nil {
			defer roomRows.Close()
			for roomRows.Next() {
				var title string
				var price float64
				var sortOrder int
				roomRows.Scan(&title, &price, &sortOrder)
				fmt.Printf("  - %s: price=%.0f\n", title, price)
			}
		}
	}

	return nil
}
