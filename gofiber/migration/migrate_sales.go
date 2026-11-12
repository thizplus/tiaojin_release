package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoSales represents sales document from MongoDB
type MongoSales struct {
	ID         primitive.ObjectID `bson:"_id"`
	Code       string             `bson:"code"`
	FirstName  string             `bson:"firstName"`
	LastName   string             `bson:"lastName"`
	NickName   string             `bson:"nickName"`
	Status     string             `bson:"status"`
	CreateDate string             `bson:"createDate"`
	UpdateDate string             `bson:"updateDate"`
}

// MigrateSales migrates sales from MongoDB to PostgreSQL
func MigrateSales(ctx context.Context, stats *MigrationStats) error {
	collection := GetMongoDB().Collection("tbl_sales")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to query MongoDB: %w", err)
	}
	defer cursor.Close(ctx)

	var salesList []MongoSales
	if err := cursor.All(ctx, &salesList); err != nil {
		return fmt.Errorf("failed to decode sales: %w", err)
	}

	fmt.Printf("    Found %d sales in MongoDB\n", len(salesList))

	skippedCount := 0

	for _, sales := range salesList {
		// Ensure code is not empty
		code := sales.Code
		if code == "" {
			code = fmt.Sprintf("SALES-%s", sales.ID.Hex()[:8])
		}

		// Check if sales already exists by code
		var existingID string
		checkQuery := `SELECT id FROM sales WHERE code = $1 AND deleted_at IS NULL LIMIT 1`
		err := GetPostgresDB().QueryRowContext(ctx, checkQuery, code).Scan(&existingID)
		if err == nil {
			// Sales already exists, set mapping and skip
			existingUUID, _ := uuid.Parse(existingID)
			SetMapping(CollectionSales, sales.ID, existingUUID)
			skippedCount++
			continue
		}

		newID := uuid.New()
		SetMapping(CollectionSales, sales.ID, newID)

		// Parse dates
		createdAt := parseDate(sales.CreateDate)
		updatedAt := parseDate(sales.UpdateDate)
		if updatedAt.IsZero() {
			updatedAt = createdAt
		}

		// Map status to active boolean
		active := mapSalesActive(sales.Status)

		// Insert into PostgreSQL
		query := `
			INSERT INTO sales (id, user_id, code, first_name, last_name, nick_name, active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`

		_, err = GetPostgresDB().ExecContext(ctx, query,
			newID,
			nil, // user_id - will be linked later if needed
			code,
			sales.FirstName,
			sales.LastName,
			sales.NickName,
			active,
			createdAt,
			updatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to insert sales %s: %w", code, err)
		}

		stats.Sales++
	}

	if skippedCount > 0 {
		fmt.Printf("    INFO: %d sales already exist (skipped)\n", skippedCount)
	}
	fmt.Printf("    Migrated %d sales\n", stats.Sales)
	return nil
}

func mapSalesActive(status string) bool {
	switch status {
	case "active", "publish", "":
		return true
	case "inactive", "draft":
		return false
	default:
		return true
	}
}
