package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FixPeriodTours fixes period.tour_id by matching tour code from MongoDB
func FixPeriodTours(ctx context.Context) error {
	fmt.Println("\n=== Fixing Period-Tour Relationships ===")

	// Step 1: Build tour code -> PostgreSQL UUID map
	fmt.Println("[1/4] Building PostgreSQL tour map...")
	pgTourMap, err := buildPGTourCodeMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to build PG tour map: %w", err)
	}
	fmt.Printf("      Found %d tours in PostgreSQL\n", len(pgTourMap))

	// Step 2: Build MongoDB tour ID -> code map
	fmt.Println("[2/4] Building MongoDB tour map...")
	mongoTourCodeMap, err := buildMongoTourIDToCodeMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to build MongoDB tour map: %w", err)
	}
	fmt.Printf("      Found %d tours in MongoDB\n", len(mongoTourCodeMap))

	// Step 3: Get periods from MongoDB
	fmt.Println("[3/4] Loading MongoDB periods...")
	mongoPeriods, err := getMongoPeriodsWithTour(ctx)
	if err != nil {
		return fmt.Errorf("failed to get MongoDB periods: %w", err)
	}
	fmt.Printf("      Found %d periods in MongoDB\n", len(mongoPeriods))

	// Step 4: Update periods in PostgreSQL
	fmt.Println("[4/4] Updating period tour_id...")

	updated := 0
	notFound := 0
	noMatch := 0

	for _, mp := range mongoPeriods {
		// Get tour code from MongoDB - TourID is a hex string
		tourCode, ok := mongoTourCodeMap[mp.TourID]
		if !ok {
			noMatch++
			continue
		}

		// Find PostgreSQL tour UUID by code
		tourUUID, ok := pgTourMap[strings.ToUpper(tourCode)]
		if !ok {
			notFound++
			continue
		}

		// Find and update period in PostgreSQL by matching start_date + end_date
		startDate := parsePeriodDate(mp.StartDate)
		endDate := parsePeriodDate(mp.EndDate)

		if startDate.IsZero() || endDate.IsZero() {
			noMatch++
			continue
		}

		err := updatePeriodTourID(ctx, tourUUID, startDate, endDate)
		if err != nil {
			continue
		}
		updated++
	}

	fmt.Printf("\n=== Fix Period-Tour Summary ===\n")
	fmt.Printf("    Updated: %d periods\n", updated)
	fmt.Printf("    Tour code not found: %d periods\n", noMatch)
	fmt.Printf("    Tour not in PG: %d periods\n", notFound)

	return nil
}

// MongoPeriodWithTour represents period with tour reference from MongoDB
type MongoPeriodWithTour struct {
	ID        primitive.ObjectID `bson:"_id"`
	TourID    string             `bson:"tour"` // This is a hex string, not ObjectID
	StartDate interface{}        `bson:"start_date"`
	EndDate   interface{}        `bson:"end_date"`
}

func getMongoPeriodsWithTour(ctx context.Context) ([]MongoPeriodWithTour, error) {
	collection := GetMongoDB().Collection("tbl_period")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var periods []MongoPeriodWithTour
	if err := cursor.All(ctx, &periods); err != nil {
		return nil, err
	}

	return periods, nil
}

func buildPGTourCodeMap(ctx context.Context) (map[string]string, error) {
	result := make(map[string]string)

	query := `SELECT id, code FROM tours WHERE deleted_at IS NULL`
	rows, err := GetPostgresDB().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, code string
		if err := rows.Scan(&id, &code); err != nil {
			continue
		}
		result[strings.ToUpper(code)] = id
	}

	return result, rows.Err()
}

func buildMongoTourIDToCodeMap(ctx context.Context) (map[string]string, error) {
	collection := GetMongoDB().Collection("tbl_tours")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]string)
	for cursor.Next(ctx) {
		var tour struct {
			ID   primitive.ObjectID `bson:"_id"`
			Code string             `bson:"code"`
		}
		if err := cursor.Decode(&tour); err != nil {
			continue
		}
		result[tour.ID.Hex()] = tour.Code
	}

	return result, cursor.Err()
}

// parsePeriodDate is defined in migrate_periods.go and handles interface{} types

func updatePeriodTourID(ctx context.Context, tourUUID string, startDate, endDate time.Time) error {
	// Convert to Thailand timezone (UTC+7) before matching
	// MongoDB stores "27 Mar midnight Thailand" as "26 Mar 17:00 UTC"
	// PostgreSQL stores the Thailand date (27 Mar)
	thailandTZ := time.FixedZone("Asia/Bangkok", 7*60*60)
	startDateTH := startDate.In(thailandTZ)
	endDateTH := endDate.In(thailandTZ)

	// Normalize to date only
	startDateOnly := time.Date(startDateTH.Year(), startDateTH.Month(), startDateTH.Day(), 0, 0, 0, 0, time.UTC)
	endDateOnly := time.Date(endDateTH.Year(), endDateTH.Month(), endDateTH.Day(), 0, 0, 0, 0, time.UTC)

	query := `
		UPDATE periods
		SET tour_id = $1, updated_at = NOW()
		WHERE start_date = $2::date AND end_date = $3::date AND deleted_at IS NULL
	`
	result, err := GetPostgresDB().ExecContext(ctx, query, tourUUID, startDateOnly, endDateOnly)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no period found")
	}

	return nil
}
