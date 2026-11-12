package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FixBookingTours fixes booking.tour_id by matching tour code from MongoDB
func FixBookingTours(ctx context.Context) error {
	fmt.Println("\n=== Fixing Booking-Tour Relationships ===")

	// Step 1: Build tour code -> PostgreSQL UUID map
	fmt.Println("[1/3] Building tour map...")
	tourMap, err := buildTourCodeMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to build tour map: %w", err)
	}
	fmt.Printf("      Found %d tours in PostgreSQL\n", len(tourMap))

	// Step 2: Build MongoDB tour ID -> code map
	mongoTourMap, err := buildMongoTourCodeMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to build MongoDB tour map: %w", err)
	}
	fmt.Printf("      Found %d tours in MongoDB\n", len(mongoTourMap))

	// Step 3: Get bookings from MongoDB
	fmt.Println("[2/3] Loading MongoDB bookings...")
	mongoBookings, err := getMongoBookingRelations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get MongoDB bookings: %w", err)
	}
	fmt.Printf("      Found %d bookings in MongoDB\n", len(mongoBookings))

	// Step 4: Update bookings in PostgreSQL
	fmt.Println("[3/3] Updating booking tour_id...")

	updated := 0
	notFound := 0

	for _, mb := range mongoBookings {
		// Get tour code from MongoDB
		tourCode, ok := mongoTourMap[mb.TourID.Hex()]
		if !ok {
			notFound++
			continue
		}

		// Find PostgreSQL tour UUID by code
		tourUUID, ok := tourMap[strings.ToUpper(tourCode)]
		if !ok {
			notFound++
			continue
		}

		// Update booking tour_id
		err := updateBookingTour(ctx, mb.InvoiceID, tourUUID)
		if err != nil {
			continue
		}
		updated++
	}

	fmt.Printf("\n=== Fix Booking-Tour Summary ===\n")
	fmt.Printf("    Updated: %d bookings\n", updated)
	fmt.Printf("    Tour not found: %d bookings\n", notFound)

	return nil
}

// FixBookingPeriods fixes booking.period_id by matching period from MongoDB
func FixBookingPeriods(ctx context.Context) error {
	fmt.Println("\n=== Fixing Booking-Period Relationships ===")

	// Step 1: Build tour code -> PostgreSQL UUID map
	tourMap, err := buildTourCodeMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to build tour map: %w", err)
	}

	// Step 2: Build MongoDB tour ID -> code map
	mongoTourMap, err := buildMongoTourCodeMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to build MongoDB tour map: %w", err)
	}

	// Step 3: Get bookings from MongoDB
	fmt.Println("[1/2] Loading MongoDB bookings...")
	mongoBookings, err := getMongoBookingRelations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get MongoDB bookings: %w", err)
	}
	fmt.Printf("      Found %d bookings in MongoDB\n", len(mongoBookings))

	// Step 4: Get MongoDB periods
	mongoPeriodMap, err := buildMongoPeriodMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to build MongoDB period map: %w", err)
	}

	// Step 5: Update bookings in PostgreSQL
	fmt.Println("[2/2] Updating booking period_id...")

	updated := 0
	notFound := 0

	for _, mb := range mongoBookings {
		// Get period info from MongoDB
		periodInfo, ok := mongoPeriodMap[mb.PeriodID.Hex()]
		if !ok {
			notFound++
			continue
		}

		// Get tour code - TourID is now a string (hex)
		tourCode, ok := mongoTourMap[periodInfo.TourID]
		if !ok {
			notFound++
			continue
		}

		// Find PostgreSQL tour UUID
		tourUUID, ok := tourMap[strings.ToUpper(tourCode)]
		if !ok {
			notFound++
			continue
		}

		// Find PostgreSQL period by tour_id + dates
		// Parse dates from interface{} to time.Time
		startDate := parsePeriodDate(periodInfo.StartDate)
		endDate := parsePeriodDate(periodInfo.EndDate)

		if startDate.IsZero() || endDate.IsZero() {
			notFound++
			continue
		}

		periodUUID, err := findPeriodByTourAndDatesTime(ctx, tourUUID, startDate, endDate)
		if err != nil || periodUUID == "" {
			notFound++
			continue
		}

		// Update booking period_id
		err = updateBookingPeriod(ctx, mb.InvoiceID, periodUUID)
		if err != nil {
			continue
		}
		updated++
	}

	fmt.Printf("\n=== Fix Booking-Period Summary ===\n")
	fmt.Printf("    Updated: %d bookings\n", updated)
	fmt.Printf("    Period not found: %d bookings\n", notFound)

	return nil
}

// FixBookingSales fixes booking.sales_id by matching sales code from MongoDB
func FixBookingSales(ctx context.Context) error {
	fmt.Println("\n=== Fixing Booking-Sales Relationships ===")

	// Step 1: Build sales code -> PostgreSQL UUID map
	fmt.Println("[1/3] Building sales map...")
	salesMap, err := buildSalesCodeMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to build sales map: %w", err)
	}
	fmt.Printf("      Found %d sales in PostgreSQL\n", len(salesMap))

	// Step 2: Build MongoDB sales ID -> code map
	mongoSalesMap, err := buildMongoSalesCodeMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to build MongoDB sales map: %w", err)
	}
	fmt.Printf("      Found %d sales in MongoDB\n", len(mongoSalesMap))

	// Step 3: Get bookings from MongoDB
	fmt.Println("[2/3] Loading MongoDB bookings...")
	mongoBookings, err := getMongoBookingRelations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get MongoDB bookings: %w", err)
	}
	fmt.Printf("      Found %d bookings in MongoDB\n", len(mongoBookings))

	// Step 4: Update bookings in PostgreSQL
	fmt.Println("[3/3] Updating booking sales_id...")

	updated := 0
	notFound := 0

	for _, mb := range mongoBookings {
		if mb.SalesID.IsZero() {
			continue
		}

		// Get sales code from MongoDB
		salesCode, ok := mongoSalesMap[mb.SalesID.Hex()]
		if !ok {
			notFound++
			continue
		}

		// Find PostgreSQL sales UUID by code
		salesUUID, ok := salesMap[strings.ToUpper(salesCode)]
		if !ok {
			notFound++
			continue
		}

		// Update booking sales_id
		err := updateBookingSales(ctx, mb.InvoiceID, salesUUID)
		if err != nil {
			continue
		}
		updated++
	}

	fmt.Printf("\n=== Fix Booking-Sales Summary ===\n")
	fmt.Printf("    Updated: %d bookings\n", updated)
	fmt.Printf("    Sales not found: %d bookings\n", notFound)

	return nil
}

// buildTourCodeMap creates map of tour code -> PostgreSQL UUID
func buildTourCodeMap(ctx context.Context) (map[string]string, error) {
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

// MongoTourCode represents tour with code from MongoDB
type MongoTourCode struct {
	ID   primitive.ObjectID `bson:"_id"`
	Code string             `bson:"code"`
}

// buildMongoTourCodeMap creates map of MongoDB tour ObjectID hex -> code
func buildMongoTourCodeMap(ctx context.Context) (map[string]string, error) {
	collection := GetMongoDB().Collection("tbl_tours")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]string)
	for cursor.Next(ctx) {
		var tour MongoTourCode
		if err := cursor.Decode(&tour); err != nil {
			continue
		}
		result[tour.ID.Hex()] = tour.Code
	}

	return result, cursor.Err()
}

// MongoBookingRelation represents booking with all relations from MongoDB
type MongoBookingRelation struct {
	ID        primitive.ObjectID `bson:"_id"`
	InvoiceID string             `bson:"invoiceId"`
	TourID    primitive.ObjectID `bson:"tour"`
	PeriodID  primitive.ObjectID `bson:"periodCode"`
	SalesID   primitive.ObjectID `bson:"salesCode"`
}

func getMongoBookingRelations(ctx context.Context) ([]MongoBookingRelation, error) {
	collection := GetMongoDB().Collection("tbl_booking")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []MongoBookingRelation
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

// MongoPeriodInfo represents period info from MongoDB
type MongoPeriodInfo struct {
	ID        primitive.ObjectID `bson:"_id"`
	TourID    string             `bson:"tour"` // This is a hex string, not ObjectID
	StartDate interface{}        `bson:"start_date"`
	EndDate   interface{}        `bson:"end_date"`
}

func buildMongoPeriodMap(ctx context.Context) (map[string]*MongoPeriodInfo, error) {
	collection := GetMongoDB().Collection("tbl_period")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]*MongoPeriodInfo)
	for cursor.Next(ctx) {
		var period MongoPeriodInfo
		if err := cursor.Decode(&period); err != nil {
			continue
		}
		result[period.ID.Hex()] = &period
	}

	return result, cursor.Err()
}

func findPeriodByTourAndDates(ctx context.Context, tourUUID, startDate, endDate string) (string, error) {
	// Parse dates
	start := parsePeriodDateStr(startDate)
	end := parsePeriodDateStr(endDate)

	if start.IsZero() || end.IsZero() {
		return "", fmt.Errorf("invalid dates")
	}

	var periodID string
	query := `
		SELECT id FROM periods
		WHERE tour_id = $1 AND start_date = $2 AND end_date = $3 AND deleted_at IS NULL
		LIMIT 1
	`
	err := GetPostgresDB().QueryRowContext(ctx, query, tourUUID, start, end).Scan(&periodID)
	if err != nil {
		return "", err
	}

	return periodID, nil
}

func findPeriodByTourAndDatesTime(ctx context.Context, tourUUID string, startDate, endDate time.Time) (string, error) {
	var periodID string
	query := `
		SELECT id FROM periods
		WHERE tour_id = $1 AND start_date = $2 AND end_date = $3 AND deleted_at IS NULL
		LIMIT 1
	`
	err := GetPostgresDB().QueryRowContext(ctx, query, tourUUID, startDate, endDate).Scan(&periodID)
	if err != nil {
		return "", err
	}

	return periodID, nil
}

func parsePeriodDateStr(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}

	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	return time.Time{}
}

func updateBookingTour(ctx context.Context, invoiceID, tourUUID string) error {
	query := `UPDATE bookings SET tour_id = $1, updated_at = NOW() WHERE invoice_id = $2 AND deleted_at IS NULL`
	_, err := GetPostgresDB().ExecContext(ctx, query, tourUUID, invoiceID)
	return err
}

func updateBookingPeriod(ctx context.Context, invoiceID, periodUUID string) error {
	query := `UPDATE bookings SET period_id = $1, updated_at = NOW() WHERE invoice_id = $2 AND deleted_at IS NULL`
	_, err := GetPostgresDB().ExecContext(ctx, query, periodUUID, invoiceID)
	return err
}

func updateBookingSales(ctx context.Context, invoiceID, salesUUID string) error {
	query := `UPDATE bookings SET sales_id = $1, updated_at = NOW() WHERE invoice_id = $2 AND deleted_at IS NULL`
	_, err := GetPostgresDB().ExecContext(ctx, query, salesUUID, invoiceID)
	return err
}

// MongoSalesCode represents sales with code from MongoDB
type MongoSalesCode struct {
	ID   primitive.ObjectID `bson:"_id"`
	Code string             `bson:"code"`
}

func buildSalesCodeMap(ctx context.Context) (map[string]string, error) {
	result := make(map[string]string)

	query := `SELECT id, code FROM sales WHERE deleted_at IS NULL`
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

func buildMongoSalesCodeMap(ctx context.Context) (map[string]string, error) {
	collection := GetMongoDB().Collection("tbl_sales")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]string)
	for cursor.Next(ctx) {
		var sales MongoSalesCode
		if err := cursor.Decode(&sales); err != nil {
			continue
		}
		result[sales.ID.Hex()] = sales.Code
	}

	return result, cursor.Err()
}
