package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoTour represents tour document from MongoDB
type MongoTour struct {
	ID        primitive.ObjectID  `bson:"_id"`
	Title     string              `bson:"title"`
	Code      string              `bson:"code"`
	CityID    *primitive.ObjectID `bson:"city"`
	CityTitle *string             `bson:"cityTitle"`
	Detail    *string             `bson:"detail"`
	Thumbnail *string             `bson:"thumbnail"`
	FilePDF   *string             `bson:"file_pdf"`
	FileWord  *string             `bson:"file_word"`
	Status    string              `bson:"status"`
	CreateDate string             `bson:"createDate"`
	UpdateDate string             `bson:"updateDate"`
}

// MigrateTours migrates tours from MongoDB to PostgreSQL
func MigrateTours(ctx context.Context, stats *MigrationStats) error {
	collection := GetMongoDB().Collection("tbl_tours")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to query MongoDB: %w", err)
	}
	defer cursor.Close(ctx)

	var tours []MongoTour
	if err := cursor.All(ctx, &tours); err != nil {
		return fmt.Errorf("failed to decode tours: %w", err)
	}

	fmt.Printf("    Found %d tours in MongoDB\n", len(tours))

	skippedCount := 0

	for _, tour := range tours {
		// Check if tour already exists by code (unique identifier)
		var existingID string
		checkQuery := `SELECT id FROM tours WHERE code = $1 AND deleted_at IS NULL LIMIT 1`
		err := GetPostgresDB().QueryRowContext(ctx, checkQuery, tour.Code).Scan(&existingID)
		if err == nil {
			// Tour already exists, set mapping and skip
			existingUUID, _ := uuid.Parse(existingID)
			SetMapping(CollectionTours, tour.ID, existingUUID)
			skippedCount++
			continue
		}

		newID := uuid.New()
		SetMapping(CollectionTours, tour.ID, newID)

		// Parse dates
		createdAt := parseDate(tour.CreateDate)
		updatedAt := parseDate(tour.UpdateDate)
		if updatedAt.IsZero() {
			updatedAt = createdAt
		}

		// Map city ID if exists
		var cityUUID *uuid.UUID
		if tour.CityID != nil {
			if mapped, ok := GetMapping(CollectionCities, *tour.CityID); ok {
				cityUUID = &mapped
			} else {
				fmt.Printf("    WARNING: City not found for tour %s (city_id: %s)\n", tour.Code, tour.CityID.Hex())
			}
		}

		// Map status
		status := mapTourStatus(tour.Status)

		// Insert into PostgreSQL
		query := `
			INSERT INTO tours (id, title, code, city_id, detail, thumbnail, file_pdf, file_word, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`

		_, err = GetPostgresDB().ExecContext(ctx, query,
			newID,
			tour.Title,
			tour.Code,
			cityUUID,
			tour.Detail,
			tour.Thumbnail,
			tour.FilePDF,
			tour.FileWord,
			status,
			createdAt,
			updatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to insert tour %s: %w", tour.Code, err)
		}

		stats.Tours++
	}

	if skippedCount > 0 {
		fmt.Printf("    INFO: %d tours already exist (skipped)\n", skippedCount)
	}
	fmt.Printf("    Migrated %d tours\n", stats.Tours)
	return nil
}

func mapTourStatus(status string) string {
	switch status {
	case "publish", "active":
		return "publish"
	case "draft", "inactive", "":
		return "draft"
	default:
		return "draft"
	}
}
