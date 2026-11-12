package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FixBookingUsers updates booking user_id by matching user from MongoDB
func FixBookingUsers(ctx context.Context) error {
	fmt.Println("\n=== Fixing Booking-User Relationships ===")

	// Step 1: Build mapping: MongoDB email -> PostgreSQL UUID
	fmt.Println("[1/4] Building user mapping...")

	// Get all users from MongoDB
	mongoColl := GetMongoDB().Collection("tbl_users")
	cursor, err := mongoColl.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to query MongoDB users: %w", err)
	}

	var mongoUsers []struct {
		ID    primitive.ObjectID `bson:"_id"`
		Email string             `bson:"email"`
	}
	if err := cursor.All(ctx, &mongoUsers); err != nil {
		cursor.Close(ctx)
		return fmt.Errorf("failed to decode MongoDB users: %w", err)
	}
	cursor.Close(ctx)

	// Build email -> mongo_id map
	emailToMongoID := make(map[string]primitive.ObjectID)
	for _, u := range mongoUsers {
		email := normalizeForLookup(u.Email)
		emailToMongoID[email] = u.ID
	}
	fmt.Printf("      MongoDB users: %d\n", len(mongoUsers))

	// Step 2: Get all users from PostgreSQL and build mapping
	fmt.Println("[2/4] Matching with PostgreSQL users...")
	rows, err := GetPostgresDB().QueryContext(ctx, "SELECT id, username FROM users WHERE deleted_at IS NULL")
	if err != nil {
		return fmt.Errorf("failed to query PostgreSQL users: %w", err)
	}
	defer rows.Close()

	// Build mongo_id -> pg_uuid map
	mongoIDToPgUUID := make(map[string]uuid.UUID)
	pgUserCount := 0
	matchCount := 0
	for rows.Next() {
		var pgID uuid.UUID
		var username string
		if err := rows.Scan(&pgID, &username); err != nil {
			continue
		}
		pgUserCount++

		// Find matching MongoDB user by username (normalized)
		normalizedUsername := normalizeForLookup(username)
		if mongoID, ok := emailToMongoID[normalizedUsername]; ok {
			mongoIDToPgUUID[mongoID.Hex()] = pgID
			matchCount++
		}
	}
	fmt.Printf("      PostgreSQL users: %d\n", pgUserCount)
	fmt.Printf("      Matched: %d\n", matchCount)

	// Step 3: Get all bookings from MongoDB
	fmt.Println("[3/4] Reading MongoDB bookings...")
	bookingColl := GetMongoDB().Collection("tbl_booking")
	bookingCursor, err := bookingColl.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to query MongoDB bookings: %w", err)
	}

	var mongoBookings []struct {
		ID        primitive.ObjectID `bson:"_id"`
		InvoiceID string             `bson:"invoiceId"`
		ProfileID primitive.ObjectID `bson:"profileId"`
	}
	if err := bookingCursor.All(ctx, &mongoBookings); err != nil {
		bookingCursor.Close(ctx)
		return fmt.Errorf("failed to decode MongoDB bookings: %w", err)
	}
	bookingCursor.Close(ctx)
	fmt.Printf("      MongoDB bookings: %d\n", len(mongoBookings))

	// Step 4: Update bookings in PostgreSQL
	fmt.Println("[4/4] Updating PostgreSQL bookings...")

	updateStmt, err := GetPostgresDB().PrepareContext(ctx, `
		UPDATE bookings SET user_id = $1, updated_at = NOW() WHERE invoice_id = $2 AND deleted_at IS NULL
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer updateStmt.Close()

	updated := 0
	notFound := 0
	noMapping := 0

	for _, b := range mongoBookings {
		// Find matching PostgreSQL user UUID
		pgUserUUID, ok := mongoIDToPgUUID[b.ProfileID.Hex()]
		if !ok {
			noMapping++
			continue
		}

		// Update booking
		result, err := updateStmt.ExecContext(ctx, pgUserUUID, b.InvoiceID)
		if err != nil {
			fmt.Printf("      ERROR: Failed to update %s: %v\n", b.InvoiceID, err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			updated++
		} else {
			notFound++
		}

		// Progress
		if updated%1000 == 0 && updated > 0 {
			fmt.Printf("      Progress: %d updated\n", updated)
		}
	}

	// Print summary
	fmt.Printf("\n=== Fix Booking Users Summary ===\n")
	fmt.Printf("    Updated: %d bookings\n", updated)
	fmt.Printf("    Not found in PG: %d bookings\n", notFound)
	fmt.Printf("    No user mapping: %d bookings\n", noMapping)

	// Verify results
	fmt.Println("\n=== Verification ===")

	var uniqueUsers int
	GetPostgresDB().QueryRowContext(ctx, "SELECT COUNT(DISTINCT user_id) FROM bookings WHERE deleted_at IS NULL").Scan(&uniqueUsers)
	fmt.Printf("    Unique user_ids in bookings: %d\n", uniqueUsers)

	var withCompany int
	GetPostgresDB().QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT b.id) FROM bookings b
		JOIN company_profiles cp ON b.user_id = cp.user_id
		WHERE b.deleted_at IS NULL
	`).Scan(&withCompany)
	fmt.Printf("    Bookings with company: %d\n", withCompany)

	return nil
}

func normalizeForLookup(s string) string {
	return strings.ToLower(strings.TrimSpace(strings.ReplaceAll(s, " ", "")))
}
