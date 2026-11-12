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

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer db.Close()

	fmt.Println("=== ตรวจสอบ Booking ใน PostgreSQL ===\n")

	// 1. ดูตัวอย่าง booking พร้อม user info
	fmt.Println("--- ตัวอย่าง Booking 5 รายการ ---")
	rows, _ := db.Query(`
		SELECT
			b.invoice_id,
			b.customer_name,
			u.username,
			u.email,
			CASE WHEN cp.id IS NOT NULL THEN 'มี' ELSE 'ไม่มี' END as has_company
		FROM bookings b
		LEFT JOIN users u ON b.user_id = u.id
		LEFT JOIN company_profiles cp ON u.id = cp.user_id
		WHERE b.deleted_at IS NULL
		ORDER BY b.created_at DESC
		LIMIT 5
	`)
	defer rows.Close()

	for rows.Next() {
		var invoiceID, customerName, username string
		var email sql.NullString
		var hasCompany string
		rows.Scan(&invoiceID, &customerName, &username, &email, &hasCompany)

		emailStr := "NULL"
		if email.Valid {
			emailStr = email.String
		}
		fmt.Printf("Invoice: %s\n", invoiceID)
		fmt.Printf("  ลูกค้า: %s\n", customerName)
		fmt.Printf("  Agent: %s\n", username)
		fmt.Printf("  Email: %s\n", emailStr)
		fmt.Printf("  บริษัท: %s\n\n", hasCompany)
	}

	// 2. สถิติ
	fmt.Println("--- สถิติ ---")

	var totalBookings int
	db.QueryRow("SELECT COUNT(*) FROM bookings WHERE deleted_at IS NULL").Scan(&totalBookings)
	fmt.Printf("Booking ทั้งหมด: %d\n", totalBookings)

	var bookingsWithCompany int
	db.QueryRow(`
		SELECT COUNT(DISTINCT b.id) FROM bookings b
		JOIN company_profiles cp ON b.user_id = cp.user_id
		WHERE b.deleted_at IS NULL
	`).Scan(&bookingsWithCompany)
	fmt.Printf("Booking ที่มีบริษัท: %d\n", bookingsWithCompany)

	var bookingsWithEmail int
	db.QueryRow(`
		SELECT COUNT(*) FROM bookings b
		JOIN users u ON b.user_id = u.id
		WHERE b.deleted_at IS NULL AND u.email IS NOT NULL
	`).Scan(&bookingsWithEmail)
	fmt.Printf("Booking ที่ Agent มี email: %d\n", bookingsWithEmail)

	// 3. เช็คว่ามี user_id ที่เป็น default user กี่รายการ
	fmt.Println("\n--- เช็ค Default User ---")
	var defaultUserID, defaultUsername string
	db.QueryRow("SELECT id, username FROM users ORDER BY created_at LIMIT 1").Scan(&defaultUserID, &defaultUsername)
	fmt.Printf("Default User: %s (%s)\n", defaultUsername, defaultUserID)

	var bookingsWithDefaultUser int
	db.QueryRow("SELECT COUNT(*) FROM bookings WHERE user_id = $1 AND deleted_at IS NULL", defaultUserID).Scan(&bookingsWithDefaultUser)
	fmt.Printf("Booking ที่ใช้ Default User: %d\n", bookingsWithDefaultUser)

	// 4. ดู user_id ที่ unique ใน bookings
	var uniqueUserIDs int
	db.QueryRow("SELECT COUNT(DISTINCT user_id) FROM bookings WHERE deleted_at IS NULL").Scan(&uniqueUserIDs)
	fmt.Printf("\nUser ID ที่ unique ใน bookings: %d\n", uniqueUserIDs)
}
