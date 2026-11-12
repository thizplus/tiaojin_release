package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	mongoDB     *mongo.Database
	postgresDB  *sql.DB
)

// InitDatabases initializes connections to both MongoDB and PostgreSQL
func InitDatabases() error {
	// Initialize MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return fmt.Errorf("MONGODB_URI environment variable is required")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping MongoDB to verify connection
	if err := client.Ping(context.Background(), nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	mongoClient = client
	mongoDB = client.Database("tiaojin")

	// Initialize PostgreSQL
	pgHost := os.Getenv("DB_HOST")
	pgPort := os.Getenv("DB_PORT")
	pgUser := os.Getenv("DB_USER")
	pgPassword := os.Getenv("DB_PASSWORD")
	pgDBName := os.Getenv("DB_NAME")
	pgSSLMode := os.Getenv("DB_SSL_MODE")

	if pgHost == "" {
		pgHost = "localhost"
	}
	if pgPort == "" {
		pgPort = "5432"
	}
	if pgSSLMode == "" {
		pgSSLMode = "disable"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pgHost, pgPort, pgUser, pgPassword, pgDBName, pgSSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	// Test PostgreSQL connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	postgresDB = db

	return nil
}

// CloseDatabases closes all database connections
func CloseDatabases() {
	if mongoClient != nil {
		mongoClient.Disconnect(context.Background())
	}
	if postgresDB != nil {
		postgresDB.Close()
	}
}

// GetMongoDB returns the MongoDB database instance
func GetMongoDB() *mongo.Database {
	return mongoDB
}

// GetPostgresDB returns the PostgreSQL database instance
func GetPostgresDB() *sql.DB {
	return postgresDB
}

// ClearPostgresData clears all data from PostgreSQL tables in correct order
func ClearPostgresData(ctx context.Context) error {
	fmt.Println("WARNING: This will delete ALL data from PostgreSQL!")
	fmt.Println()

	// Order matters due to foreign key constraints - delete children first
	tables := []string{
		// Sync tracking
		"sync_tracking",
		// Booking related (children first)
		"booking_additional_services",
		"passenger_passports",
		"payment_slips",
		"booking_rate_types",
		"booking_room_types",
		"bookings",
		// Period related (children first)
		"period_additional_services",
		"period_room_types",
		"period_rate_items",
		"periods",
		// Tours
		"tours",
		// Independent tables
		"sales",
		"files",           // files table has FK to users
		"api_keys",        // api_keys table has FK to users
		"seat_holds",      // seat_holds may have FK to users
		"company_profiles",
		"users",
		"banners",
		"cities",
	}

	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM %s", table)
		if _, err := postgresDB.ExecContext(ctx, query); err != nil {
			// Ignore errors for tables that might not exist
			fmt.Printf("      - Skipped table (may not exist): %s\n", table)
			continue
		}
		fmt.Printf("      - Cleared table: %s\n", table)
	}

	// Also clear any ID sequences/reset auto-increment if needed
	fmt.Println()
	fmt.Println("All tables cleared successfully!")

	return nil
}
