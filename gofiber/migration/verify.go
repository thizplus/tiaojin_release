package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// VerificationResult holds the count comparison for a collection/table
type VerificationResult struct {
	Name         string
	MongoCount   int64
	PostgresCount int64
	Difference   int64
	Status       string
}

// VerifyMigration compares counts between MongoDB collections and PostgreSQL tables
func VerifyMigration(ctx context.Context) error {
	fmt.Println()
	fmt.Println("============================================")
	fmt.Println("  Migration Verification Report")
	fmt.Println("============================================")
	fmt.Println()

	results := []VerificationResult{}

	// Define collection/table mappings
	verifications := []struct {
		name          string
		mongoCollection string
		postgresTable string
	}{
		{"Cities", "tbl_city", "cities"},
		{"Banners", "tbl_banner", "banners"},
		{"Sales", "tbl_sales", "sales"},
		{"Users", "tbl_users", "users"},
		{"Company Profiles", "tbl_users", "company_profiles"},
		{"Tours", "tbl_tours", "tours"},
		{"Periods", "tbl_period", "periods"},
		{"Bookings", "tbl_booking", "bookings"},
	}

	hasIssues := false

	for _, v := range verifications {
		mongoCount, err := getMongoCount(ctx, v.mongoCollection)
		if err != nil {
			fmt.Printf("  ERROR: Failed to count MongoDB %s: %v\n", v.mongoCollection, err)
			continue
		}

		postgresCount, err := getPostgresCount(ctx, v.postgresTable)
		if err != nil {
			fmt.Printf("  ERROR: Failed to count PostgreSQL %s: %v\n", v.postgresTable, err)
			continue
		}

		diff := mongoCount - postgresCount
		status := "OK"
		if diff > 0 {
			status = "MISSING"
			hasIssues = true
		} else if diff < 0 {
			status = "EXTRA"
			hasIssues = true
		}

		results = append(results, VerificationResult{
			Name:         v.name,
			MongoCount:   mongoCount,
			PostgresCount: postgresCount,
			Difference:   diff,
			Status:       status,
		})
	}

	// Print results in a table format
	fmt.Printf("  %-20s %10s %10s %10s %10s\n", "Entity", "MongoDB", "Postgres", "Diff", "Status")
	fmt.Println("  " + "------------------------------------------------------------")

	for _, r := range results {
		statusIcon := "OK"
		if r.Status == "MISSING" {
			statusIcon = "MISSING"
		} else if r.Status == "EXTRA" {
			statusIcon = "EXTRA"
		}
		fmt.Printf("  %-20s %10d %10d %10d %10s\n", r.Name, r.MongoCount, r.PostgresCount, r.Difference, statusIcon)
	}

	fmt.Println()

	// Additional child table verifications
	fmt.Println("  Child Tables (PostgreSQL only):")
	fmt.Println("  " + "------------------------------------------------------------")

	childTables := []string{
		"period_rate_items",
		"period_room_types",
		"booking_room_types",
		"booking_rate_types",
		"payment_slips",
		"passenger_passports",
		"booking_additional_services",
	}

	for _, table := range childTables {
		count, err := getPostgresCount(ctx, table)
		if err != nil {
			fmt.Printf("  ERROR: Failed to count %s: %v\n", table, err)
			continue
		}
		fmt.Printf("  %-30s %10d records\n", table, count)
	}

	fmt.Println()

	if hasIssues {
		fmt.Println("  WARNING: Some discrepancies found. Please review the migration.")
		return fmt.Errorf("verification found discrepancies")
	}

	fmt.Println("  All counts match! Migration verified successfully.")
	return nil
}

func getMongoCount(ctx context.Context, collection string) (int64, error) {
	coll := GetMongoDB().Collection(collection)
	return coll.CountDocuments(ctx, bson.M{})
}

func getPostgresCount(ctx context.Context, table string) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE deleted_at IS NULL", table)
	err := GetPostgresDB().QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		// Try without deleted_at filter (for tables without soft delete)
		query = fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
		err = GetPostgresDB().QueryRowContext(ctx, query).Scan(&count)
	}
	return count, err
}
