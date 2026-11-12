package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CheckTourPeriods checks periods for a specific tour code
func CheckTourPeriods(ctx context.Context, tourCode string) error {
	fmt.Printf("\n=== Checking Periods for Tour: %s ===\n\n", tourCode)

	// Check PostgreSQL
	fmt.Println("--- PostgreSQL Periods ---")
	rows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT p.id, p.start_date, p.end_date, p.seat, p.status, p.created_at
		FROM periods p
		JOIN tours t ON p.tour_id = t.id
		WHERE t.code = $1 AND p.deleted_at IS NULL
		ORDER BY p.start_date, p.created_at
	`, tourCode)
	if err != nil {
		return err
	}
	defer rows.Close()

	pgCount := 0
	for rows.Next() {
		var id, status string
		var startDate, endDate, createdAt interface{}
		var seat int
		if err := rows.Scan(&id, &startDate, &endDate, &seat, &status, &createdAt); err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		pgCount++
		fmt.Printf("%d. ID: %s\n", pgCount, id[:8]+"...")
		fmt.Printf("   Start: %v | End: %v | Seat: %d | Status: %s\n", startDate, endDate, seat, status)
		fmt.Printf("   Created: %v\n\n", createdAt)
	}
	fmt.Printf("Total in PostgreSQL: %d\n\n", pgCount)

	// Check MongoDB
	fmt.Println("--- MongoDB Periods ---")

	// First get ALL tours with this code from MongoDB
	tourCollection := GetMongoDB().Collection("tbl_tours")

	// Count all tours with this code
	tourCursor, err := tourCollection.Find(ctx, bson.M{"code": tourCode})
	if err != nil {
		return fmt.Errorf("tour query failed: %v", err)
	}

	fmt.Println("--- All MongoDB Tours with this code ---")
	var allTourIDs []string
	tourCount := 0
	for tourCursor.Next(ctx) {
		var t map[string]interface{}
		if err := tourCursor.Decode(&t); err != nil {
			continue
		}
		tourCount++
		tourID := t["_id"].(primitive.ObjectID).Hex()
		allTourIDs = append(allTourIDs, tourID)
		fmt.Printf("  Tour %d: _id=%s, code=%v, createDate=%v\n", tourCount, tourID, t["code"], t["createDate"])
	}
	tourCursor.Close(ctx)
	fmt.Printf("Total tours with code %s: %d\n\n", tourCode, tourCount)

	if tourCount == 0 {
		return fmt.Errorf("tour not found in MongoDB")
	}

	// Use first tour for compatibility
	var tour struct {
		ID   interface{} `bson:"_id"`
		Code string      `bson:"code"`
	}
	err = tourCollection.FindOne(ctx, bson.M{"code": tourCode}).Decode(&tour)
	if err != nil {
		return fmt.Errorf("tour not found in MongoDB: %v", err)
	}
	fmt.Printf("Using MongoDB Tour ID: %v\n\n", tour.ID)

	// Check periods for ALL tour IDs with this code
	fmt.Println("--- MongoDB Periods for ALL tour IDs ---")
	totalPeriodsAllTours := 0
	periodCollection := GetMongoDB().Collection("tbl_period")
	for _, tid := range allTourIDs {
		cursor, _ := periodCollection.Find(ctx, bson.M{"tour": tid})
		count := 0
		for cursor.Next(ctx) {
			var period map[string]interface{}
			if err := cursor.Decode(&period); err != nil {
				continue
			}
			count++
			totalPeriodsAllTours++
			fmt.Printf("  Tour %s -> Period: _id=%v, start=%v, createDate=%v\n",
				tid[:8]+"...", period["_id"], period["start_date"], period["createDate"])
		}
		cursor.Close(ctx)
		if count == 0 {
			fmt.Printf("  Tour %s -> No periods\n", tid[:8]+"...")
		}
	}
	fmt.Printf("Total periods for all tour IDs: %d\n\n", totalPeriodsAllTours)

	// Get periods for this tour
	// MongoDB uses field "tour" (not "tour_id") and stores as hex string
	// periodCollection already declared above

	// Get tour ObjectID hex string (without "ObjectID(" prefix)
	tourHex := tour.ID.(primitive.ObjectID).Hex()
	fmt.Printf("Searching periods with tour = %s\n\n", tourHex)

	cursor, err := periodCollection.Find(ctx, bson.M{"tour": tourHex})
	if err != nil {
		return err
	}

	mongoCount := 0
	for cursor.Next(ctx) {
		var period map[string]interface{}
		if err := cursor.Decode(&period); err != nil {
			continue
		}
		mongoCount++
		fmt.Printf("%d. MongoDB ID: %v\n", mongoCount, period["_id"])
		fmt.Printf("   tour: %v\n", period["tour"])
		fmt.Printf("   start_date: %v\n", period["start_date"])
		fmt.Printf("   end_date: %v\n", period["end_date"])
		fmt.Printf("   seat: %v | status: %v\n\n", period["seat"], period["status"])
	}
	cursor.Close(ctx)

	fmt.Printf("Total in MongoDB: %d\n", mongoCount)

	// Also check for ALL periods containing this tour ID anywhere
	fmt.Println("\n--- Searching ALL periods containing tour ID ---")
	allCursor, _ := periodCollection.Find(ctx, bson.M{})
	matchCount := 0
	for allCursor.Next(ctx) {
		var period map[string]interface{}
		if err := allCursor.Decode(&period); err != nil {
			continue
		}
		tourVal := fmt.Sprintf("%v", period["tour"])
		if tourVal == tourHex {
			matchCount++
			fmt.Printf("  Match %d: _id=%v, createDate=%v, start_date=%v\n",
				matchCount, period["_id"], period["createDate"], period["start_date"])
		}
	}
	allCursor.Close(ctx)
	fmt.Printf("Total matching periods in ALL MongoDB: %d\n", matchCount)

	// Also check ALL periods in MongoDB to see total
	totalCursor, _ := periodCollection.Find(ctx, bson.M{})
	totalCount := 0
	for totalCursor.Next(ctx) {
		totalCount++
	}
	totalCursor.Close(ctx)
	fmt.Printf("Total periods in MongoDB (all tours): %d\n", totalCount)

	// Debug: Show first 3 periods from MongoDB to see ALL field names
	fmt.Println("\n--- Debug: Sample MongoDB Periods (first 3) - ALL FIELDS ---")
	sampleCursor, _ := periodCollection.Find(ctx, bson.M{})
	sampleCount := 0
	for sampleCursor.Next(ctx) && sampleCount < 3 {
		var period map[string]interface{}
		if err := sampleCursor.Decode(&period); err != nil {
			continue
		}
		sampleCount++
		fmt.Printf("%d. Fields: ", sampleCount)
		for k, v := range period {
			fmt.Printf("%s(%T)=%v, ", k, v, v)
		}
		fmt.Println()
	}
	sampleCursor.Close(ctx)

	// Try to find periods by searching tour_id as string hex
	fmt.Println("\n--- Searching by tour_id as hex string ---")
	hexStr := fmt.Sprintf("%x", tour.ID)
	cursor3, _ := periodCollection.Find(ctx, bson.M{"tour_id": hexStr})
	hexCount := 0
	for cursor3.Next(ctx) {
		hexCount++
	}
	cursor3.Close(ctx)
	fmt.Printf("Found %d periods with tour_id as hex string\n", hexCount)

	// Summary
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("MongoDB: %d periods\n", mongoCount)
	fmt.Printf("PostgreSQL: %d periods\n", pgCount)
	if pgCount > mongoCount {
		fmt.Printf("⚠️  DUPLICATES: %d extra periods in PostgreSQL!\n", pgCount-mongoCount)
	}

	return nil
}
