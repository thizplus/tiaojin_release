//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	godotenv.Load()
	
	ctx := context.Background()
	
	// Connect to MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	
	db := client.Database("tiaojin")
	
	fmt.Println("=== Checking Orphan Records Dates ===\n")
	
	// Get all tour IDs
	tourIDs := make(map[string]bool)
	tourCursor, _ := db.Collection("tbl_tour").Find(ctx, bson.M{})
	var tours []bson.M
	tourCursor.All(ctx, &tours)
	for _, t := range tours {
		if id, ok := t["_id"].(primitive.ObjectID); ok {
			tourIDs[id.Hex()] = true
		}
	}
	
	// Check orphan periods
	fmt.Println("--- Orphan Periods (missing tour) ---")
	periodCursor, _ := db.Collection("tbl_period").Find(ctx, bson.M{})
	var periods []bson.M
	periodCursor.All(ctx, &periods)
	
	for _, p := range periods {
		tourID := ""
		if t, ok := p["tour"].(string); ok {
			tourID = t
		}
		
		if tourID == "" || !tourIDs[tourID] {
			createDate := "unknown"
			if cd, ok := p["createDate"].(string); ok {
				createDate = cd
			}
			periodID := p["_id"].(primitive.ObjectID).Hex()
			fmt.Printf("  Period: %s | Tour: %s | Created: %s\n", periodID[:12], tourID[:min(12, len(tourID))], createDate)
		}
	}
	
	// Get all period IDs
	periodIDs := make(map[string]bool)
	for _, p := range periods {
		if id, ok := p["_id"].(primitive.ObjectID); ok {
			periodIDs[id.Hex()] = true
		}
	}
	
	// Check orphan bookings
	fmt.Println("\n--- Orphan Bookings (missing period) ---")
	bookingCursor, _ := db.Collection("tbl_booking").Find(ctx, bson.M{})
	var bookings []bson.M
	bookingCursor.All(ctx, &bookings)
	
	for _, b := range bookings {
		periodID := ""
		if p, ok := b["period"].(string); ok {
			periodID = p
		}
		
		if periodID == "" || !periodIDs[periodID] {
			createDate := "unknown"
			if cd, ok := b["createDate"].(string); ok {
				createDate = cd
			}
			invoiceID := ""
			if inv, ok := b["invoiceId"].(string); ok {
				invoiceID = inv
			}
			fmt.Printf("  Booking: %s | Period: %s | Created: %s\n", invoiceID, periodID[:min(12, len(periodID))], createDate)
		}
	}
	
	fmt.Println("\n=== Done ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
