package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoBanner represents banner document from MongoDB
type MongoBanner struct {
	ID         primitive.ObjectID `bson:"_id"`
	Image      string             `bson:"img"`
	Status     string             `bson:"status"`
	CreateDate string             `bson:"createDate"`
	UpdateDate string             `bson:"updateDate"`
}

// MigrateBanners migrates banners from MongoDB to PostgreSQL
func MigrateBanners(ctx context.Context, stats *MigrationStats) error {
	collection := GetMongoDB().Collection("tbl_banner")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to query MongoDB: %w", err)
	}
	defer cursor.Close(ctx)

	var banners []MongoBanner
	if err := cursor.All(ctx, &banners); err != nil {
		return fmt.Errorf("failed to decode banners: %w", err)
	}

	fmt.Printf("    Found %d banners in MongoDB\n", len(banners))

	skippedCount := 0

	// Get current max display_order
	var maxOrder int
	GetPostgresDB().QueryRowContext(ctx, `SELECT COALESCE(MAX(display_order), 0) FROM banners WHERE deleted_at IS NULL`).Scan(&maxOrder)

	for _, banner := range banners {
		// Check if banner already exists by image URL
		var existingID string
		checkQuery := `SELECT id FROM banners WHERE img = $1 AND deleted_at IS NULL LIMIT 1`
		err := GetPostgresDB().QueryRowContext(ctx, checkQuery, banner.Image).Scan(&existingID)
		if err == nil {
			// Banner already exists, set mapping and skip
			existingUUID, _ := uuid.Parse(existingID)
			SetMapping(CollectionBanners, banner.ID, existingUUID)
			skippedCount++
			continue
		}

		newID := uuid.New()
		SetMapping(CollectionBanners, banner.ID, newID)

		// Parse dates
		createdAt := parseDate(banner.CreateDate)
		updatedAt := parseDate(banner.UpdateDate)
		if updatedAt.IsZero() {
			updatedAt = createdAt
		}

		// Map status
		status := mapBannerStatus(banner.Status)

		// Create a title from the image filename or use default
		maxOrder++
		title := fmt.Sprintf("Banner %d", maxOrder)

		// Insert into PostgreSQL
		query := `
			INSERT INTO banners (id, title, img, link_url, alt, display_order, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`

		_, err = GetPostgresDB().ExecContext(ctx, query,
			newID,
			title,
			banner.Image,
			nil,      // link_url - not in old system
			nil,      // alt - not in old system
			maxOrder, // display_order
			status,
			createdAt,
			updatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to insert banner %s: %w", banner.ID.Hex(), err)
		}

		stats.Banners++
	}

	if skippedCount > 0 {
		fmt.Printf("    INFO: %d banners already exist (skipped)\n", skippedCount)
	}
	fmt.Printf("    Migrated %d banners\n", stats.Banners)
	return nil
}

func mapBannerStatus(status string) string {
	switch status {
	case "publish", "active":
		return "publish"
	case "draft", "inactive", "":
		return "draft"
	default:
		return "draft"
	}
}
