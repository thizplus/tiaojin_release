//go:build ignore

package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	// MongoDB
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		fmt.Println("MongoDB error:", err)
		return
	}
	defer mongoClient.Disconnect(ctx)
	mongoDB := mongoClient.Database("tiaojin")

	// PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))
	pgDB, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("PostgreSQL error:", err)
		return
	}
	defer pgDB.Close()

	fmt.Println("=== MongoDB Sample URLs ===\n")

	// MongoDB banners
	fmt.Println("--- MongoDB tbl_banner.img ---")
	cursor, _ := mongoDB.Collection("tbl_banner").Find(ctx, bson.M{}, options.Find().SetLimit(3))
	var banners []bson.M
	cursor.All(ctx, &banners)
	for _, b := range banners {
		if img, ok := b["img"].(string); ok && img != "" {
			fmt.Printf("  %s\n", truncate(img, 80))
		} else {
			fmt.Println("  (empty)")
		}
	}

	// MongoDB cities
	fmt.Println("\n--- MongoDB tbl_city.img ---")
	cursor, _ = mongoDB.Collection("tbl_city").Find(ctx, bson.M{}, options.Find().SetLimit(3))
	var cities []bson.M
	cursor.All(ctx, &cities)
	for _, c := range cities {
		if img, ok := c["img"].(string); ok && img != "" {
			fmt.Printf("  %s\n", truncate(img, 80))
		} else {
			fmt.Println("  (empty)")
		}
	}

	fmt.Println("\n=== PostgreSQL Sample URLs ===\n")

	// PostgreSQL banners
	fmt.Println("--- PostgreSQL banners.img ---")
	rows, _ := pgDB.Query("SELECT COALESCE(img, '(NULL)') FROM banners LIMIT 3")
	for rows.Next() {
		var img string
		rows.Scan(&img)
		fmt.Printf("  %s\n", truncate(img, 80))
	}
	rows.Close()

	// PostgreSQL cities
	fmt.Println("\n--- PostgreSQL cities.img ---")
	rows, _ = pgDB.Query("SELECT COALESCE(img, '(NULL)') FROM cities LIMIT 3")
	for rows.Next() {
		var img string
		rows.Scan(&img)
		fmt.Printf("  %s\n", truncate(img, 80))
	}
	rows.Close()

	// Count non-empty
	fmt.Println("\n=== PostgreSQL Non-Empty Counts ===")
	var count int
	pgDB.QueryRow("SELECT COUNT(*) FROM banners WHERE img IS NOT NULL AND img != ''").Scan(&count)
	fmt.Printf("  banners with img: %d\n", count)
	pgDB.QueryRow("SELECT COUNT(*) FROM cities WHERE img IS NOT NULL AND img != ''").Scan(&count)
	fmt.Printf("  cities with img: %d\n", count)
	pgDB.QueryRow("SELECT COUNT(*) FROM company_profiles WHERE logo IS NOT NULL AND logo != ''").Scan(&count)
	fmt.Printf("  company_profiles with logo: %d\n", count)
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
