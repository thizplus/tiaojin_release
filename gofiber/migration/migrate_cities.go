package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoCity represents city document from MongoDB
type MongoCity struct {
	ID         primitive.ObjectID `bson:"_id"`
	Title      string             `bson:"title"`
	Detail     *string            `bson:"detail"`
	Image      string             `bson:"img"`
	Status     string             `bson:"status"`
	CreateDate string             `bson:"createDate"`
	UpdateDate string             `bson:"updateDate"`
}

// MigrateCities migrates cities from MongoDB to PostgreSQL
func MigrateCities(ctx context.Context, stats *MigrationStats) error {
	collection := GetMongoDB().Collection("tbl_city")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to query MongoDB: %w", err)
	}
	defer cursor.Close(ctx)

	var cities []MongoCity
	if err := cursor.All(ctx, &cities); err != nil {
		return fmt.Errorf("failed to decode cities: %w", err)
	}

	fmt.Printf("    Found %d cities in MongoDB\n", len(cities))

	skippedCount := 0

	for _, city := range cities {
		// Check if city already exists by title (unique identifier)
		var existingID string
		checkQuery := `SELECT id FROM cities WHERE title = $1 AND deleted_at IS NULL LIMIT 1`
		err := GetPostgresDB().QueryRowContext(ctx, checkQuery, city.Title).Scan(&existingID)
		if err == nil {
			// City already exists, set mapping and skip
			existingUUID, _ := uuid.Parse(existingID)
			SetMapping(CollectionCities, city.ID, existingUUID)
			skippedCount++
			continue
		}

		newID := uuid.New()
		SetMapping(CollectionCities, city.ID, newID)

		// Parse dates
		createdAt := parseDate(city.CreateDate)
		updatedAt := parseDate(city.UpdateDate)
		if updatedAt.IsZero() {
			updatedAt = createdAt
		}

		// Map status
		status := mapCityStatus(city.Status)

		// Insert into PostgreSQL
		query := `
			INSERT INTO cities (id, title, detail, img, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`

		_, err = GetPostgresDB().ExecContext(ctx, query,
			newID,
			city.Title,
			city.Detail,
			city.Image,
			status,
			createdAt,
			updatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to insert city %s: %w", city.Title, err)
		}

		stats.Cities++
	}

	if skippedCount > 0 {
		fmt.Printf("    INFO: %d cities already exist (skipped)\n", skippedCount)
	}
	fmt.Printf("    Migrated %d cities\n", stats.Cities)
	return nil
}

func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}

	// Try various date formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	return time.Now()
}

func mapCityStatus(status string) string {
	switch status {
	case "publish", "active":
		return "publish"
	case "draft", "inactive", "":
		return "draft"
	default:
		return "draft"
	}
}
