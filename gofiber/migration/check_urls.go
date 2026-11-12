//go:build ignore

package main

import (
	"database/sql"
	"fmt"
	"os"
	
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	
	db, _ := sql.Open("postgres", connStr)
	defer db.Close()
	
	fmt.Println("=== Sample file URLs in PostgreSQL ===\n")
	
	// Company profiles
	fmt.Println("--- company_profiles.logo (first 3) ---")
	rows, _ := db.Query("SELECT logo FROM company_profiles WHERE logo IS NOT NULL AND logo != '' LIMIT 3")
	for rows.Next() {
		var url string
		rows.Scan(&url)
		if len(url) > 80 {
			fmt.Printf("  %s...\n", url[:80])
		} else {
			fmt.Printf("  %s\n", url)
		}
	}
	rows.Close()
	
	// Banners
	fmt.Println("\n--- banners.img (first 3) ---")
	rows, _ = db.Query("SELECT img FROM banners WHERE img IS NOT NULL AND img != '' LIMIT 3")
	for rows.Next() {
		var url string
		rows.Scan(&url)
		if len(url) > 80 {
			fmt.Printf("  %s...\n", url[:80])
		} else {
			fmt.Printf("  %s\n", url)
		}
	}
	rows.Close()
	
	// Tours
	fmt.Println("\n--- tours.thumbnail (first 3) ---")
	rows, _ = db.Query("SELECT thumbnail FROM tours WHERE thumbnail IS NOT NULL AND thumbnail != '' LIMIT 3")
	for rows.Next() {
		var url string
		rows.Scan(&url)
		if len(url) > 80 {
			fmt.Printf("  %s...\n", url[:80])
		} else {
			fmt.Printf("  %s\n", url)
		}
	}
	rows.Close()
	
	// Count by URL type
	fmt.Println("\n=== URL Type Counts ===")
	
	var doCount, e2Count, pathCount, emptyCount int
	db.QueryRow("SELECT COUNT(*) FROM company_profiles WHERE logo LIKE 'https://tiaojin.sgp1%'").Scan(&doCount)
	db.QueryRow("SELECT COUNT(*) FROM company_profiles WHERE logo LIKE 'https://tiaojin.s3.ap-southeast-1.idrivee2%'").Scan(&e2Count)
	db.QueryRow("SELECT COUNT(*) FROM company_profiles WHERE logo LIKE 'uploads/%'").Scan(&pathCount)
	db.QueryRow("SELECT COUNT(*) FROM company_profiles WHERE logo IS NULL OR logo = ''").Scan(&emptyCount)
	
	fmt.Printf("\ncompany_profiles.logo:\n")
	fmt.Printf("  DO Spaces URL: %d\n", doCount)
	fmt.Printf("  iDrive e2 URL: %d\n", e2Count)
	fmt.Printf("  Path only:     %d\n", pathCount)
	fmt.Printf("  Empty/NULL:    %d\n", emptyCount)
}
