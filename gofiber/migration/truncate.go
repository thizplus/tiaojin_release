//go:build ignore
// +build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Tables in order (respecting foreign keys - child tables first)
	tables := []string{
		"payment_slips",
		"passengers_passports",
		"booking_rate_types",
		"booking_room_types",
		"bookings",
		"periods",
		"tours",
		"company_profiles",
		"users",
		"sales",
		"banners",
		"cities",
	}

	fmt.Println("Truncating tables...")
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			fmt.Printf("  WARNING: Failed to truncate %s: %v\n", table, err)
		} else {
			fmt.Printf("  Truncated: %s\n", table)
		}
	}
	fmt.Println("Done!")
}
