package main

import (
	"context"
	"database/sql"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FixTourCities updates tour city_id by matching city from MongoDB
func FixTourCities(ctx context.Context) error {
	fmt.Println("\n=== Fixing Tour City Relationships ===")

	// Step 1: Build city title -> PostgreSQL UUID map
	cityMap, err := buildCityMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to build city map: %w", err)
	}
	fmt.Printf("    Found %d cities in PostgreSQL\n", len(cityMap))

	// Step 2: Get all tours from MongoDB with their city info
	mongoTours, err := getMongoToursWithCity(ctx)
	if err != nil {
		return fmt.Errorf("failed to get MongoDB tours: %w", err)
	}
	fmt.Printf("    Found %d tours in MongoDB\n", len(mongoTours))

	// Step 3: Update each tour in PostgreSQL
	updatedCount := 0
	skippedCount := 0
	notFoundCount := 0

	for _, mt := range mongoTours {
		// Skip if no city info
		if mt.CityTitle == nil || *mt.CityTitle == "" {
			skippedCount++
			continue
		}

		// Find city UUID by title
		cityUUID, ok := cityMap[*mt.CityTitle]
		if !ok {
			fmt.Printf("    WARNING: City '%s' not found for tour %s\n", *mt.CityTitle, mt.Code)
			notFoundCount++
			continue
		}

		// Update tour in PostgreSQL
		err := updateTourCity(ctx, mt.Code, cityUUID)
		if err != nil {
			fmt.Printf("    ERROR: Failed to update tour %s: %v\n", mt.Code, err)
			continue
		}
		updatedCount++
	}

	fmt.Printf("\n=== Fix Tour Cities Summary ===\n")
	fmt.Printf("    Updated: %d tours\n", updatedCount)
	fmt.Printf("    Skipped (no city): %d tours\n", skippedCount)
	fmt.Printf("    City not found: %d tours\n", notFoundCount)

	return nil
}

// buildCityMap creates a map of city title -> PostgreSQL UUID
func buildCityMap(ctx context.Context) (map[string]string, error) {
	cityMap := make(map[string]string)

	query := `SELECT id, title FROM cities WHERE deleted_at IS NULL`
	rows, err := GetPostgresDB().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, title string
		if err := rows.Scan(&id, &title); err != nil {
			return nil, err
		}
		cityMap[title] = id
	}

	return cityMap, rows.Err()
}

// MongoTourCity represents tour with city info from MongoDB
type MongoTourCity struct {
	ID        primitive.ObjectID  `bson:"_id"`
	Code      string              `bson:"code"`
	CityID    *primitive.ObjectID `bson:"city"`
	CityTitle *string             `bson:"cityTitle"`
}

// getMongoToursWithCity gets all tours from MongoDB with city information
func getMongoToursWithCity(ctx context.Context) ([]MongoTourCity, error) {
	collection := GetMongoDB().Collection("tbl_tours")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tours []MongoTourCity
	if err := cursor.All(ctx, &tours); err != nil {
		return nil, err
	}

	return tours, nil
}

// updateTourCity updates a tour's city_id in PostgreSQL
func updateTourCity(ctx context.Context, tourCode string, cityUUID string) error {
	query := `UPDATE tours SET city_id = $1, updated_at = NOW() WHERE code = $2 AND deleted_at IS NULL`
	result, err := GetPostgresDB().ExecContext(ctx, query, cityUUID, tourCode)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
