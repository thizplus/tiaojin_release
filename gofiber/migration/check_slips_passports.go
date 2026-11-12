package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// CheckSlipsAndPassports checks slip and passport data in MongoDB
func CheckSlipsAndPassports(ctx context.Context) error {
	fmt.Println("\n=== Checking Slips and Passports in MongoDB ===\n")

	collection := GetMongoDB().Collection("tbl_booking")

	// Count total bookings
	totalBookings, _ := collection.CountDocuments(ctx, bson.M{})
	fmt.Printf("Total bookings in MongoDB: %d\n\n", totalBookings)

	// Count bookings with slips
	bookingsWithSlips, _ := collection.CountDocuments(ctx, bson.M{
		"slip": bson.M{"$exists": true, "$ne": nil, "$not": bson.M{"$size": 0}},
	})
	fmt.Printf("Bookings with slip array (non-empty): %d\n", bookingsWithSlips)

	// Count bookings with passports
	bookingsWithPassports, _ := collection.CountDocuments(ctx, bson.M{
		"passport": bson.M{"$exists": true, "$ne": nil, "$not": bson.M{"$size": 0}},
	})
	fmt.Printf("Bookings with passport array (non-empty): %d\n", bookingsWithPassports)

	// Sample a few bookings with slips - show ALL fields
	fmt.Println("\n--- Sample Bookings with Slips (Full Structure) ---")
	cursor, err := collection.Find(ctx, bson.M{
		"slip": bson.M{"$exists": true, "$ne": nil, "$not": bson.M{"$size": 0}},
	})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) && count < 3 {
		var booking bson.M
		if err := cursor.Decode(&booking); err != nil {
			continue
		}

		invoiceID := booking["invoiceId"]
		sumPayment := booking["sumPayment"]
		slips := booking["slip"]

		fmt.Printf("\nBooking: %v (sumPayment: %v)\n", invoiceID, sumPayment)
		if slipArr, ok := slips.(bson.A); ok {
			fmt.Printf("  Slips count: %d\n", len(slipArr))
			for i, s := range slipArr {
				if i >= 2 {
					fmt.Printf("  ... and %d more\n", len(slipArr)-2)
					break
				}
				if slip, ok := s.(bson.M); ok {
					fmt.Printf("  Slip %d (ALL FIELDS):\n", i+1)
					for k, v := range slip {
						fmt.Printf("    %s: %v\n", k, v)
					}
				}
			}
		}
		count++
	}

	// Sample a few bookings with passports - show ALL fields
	fmt.Println("\n--- Sample Bookings with Passports (Full Structure) ---")
	cursor2, err := collection.Find(ctx, bson.M{
		"passport": bson.M{"$exists": true, "$ne": nil, "$not": bson.M{"$size": 0}},
	})
	if err != nil {
		return err
	}
	defer cursor2.Close(ctx)

	count = 0
	for cursor2.Next(ctx) && count < 2 {
		var booking bson.M
		if err := cursor2.Decode(&booking); err != nil {
			continue
		}

		invoiceID := booking["invoiceId"]
		passports := booking["passport"]

		fmt.Printf("\nBooking: %v\n", invoiceID)
		if passArr, ok := passports.(bson.A); ok {
			fmt.Printf("  Passports count: %d\n", len(passArr))
			for i, p := range passArr {
				if i >= 1 {
					fmt.Printf("  ... and %d more\n", len(passArr)-1)
					break
				}
				if passport, ok := p.(bson.M); ok {
					fmt.Printf("  Passport %d (ALL FIELDS):\n", i+1)
					for k, v := range passport {
						if k == "personal" {
							fmt.Printf("    %s:\n", k)
							if personal, ok := v.(bson.M); ok {
								for pk, pv := range personal {
									fmt.Printf("      %s: %v\n", pk, pv)
								}
							}
						} else {
							fmt.Printf("    %s: %v\n", k, v)
						}
					}
				}
			}
		}
		count++
	}

	// Check confirmed slips
	fmt.Println("\n--- Analyzing Slips with sumPayment > 0 ---")
	cursor3, _ := collection.Find(ctx, bson.M{
		"slip":       bson.M{"$exists": true, "$ne": nil, "$not": bson.M{"$size": 0}},
		"sumPayment": bson.M{"$gt": 0},
	})
	defer cursor3.Close(ctx)

	totalWithPayment := 0
	confirmedSlips := 0
	for cursor3.Next(ctx) {
		var booking bson.M
		if err := cursor3.Decode(&booking); err != nil {
			continue
		}
		totalWithPayment++

		if slipArr, ok := booking["slip"].(bson.A); ok {
			for _, s := range slipArr {
				if slip, ok := s.(bson.M); ok {
					status := slip["status"]
					if status == "confirm" || status == "confirmed" || status == "approved" {
						confirmedSlips++
					}
				}
			}
		}
	}
	fmt.Printf("Bookings with slips AND sumPayment > 0: %d\n", totalWithPayment)
	fmt.Printf("Total confirmed/approved slips: %d\n", confirmedSlips)

	// Check PostgreSQL counts
	fmt.Println("\n--- PostgreSQL Counts ---")

	var pgSlipsCount int
	err = GetPostgresDB().QueryRowContext(ctx, "SELECT COUNT(*) FROM payment_slips WHERE deleted_at IS NULL").Scan(&pgSlipsCount)
	if err != nil {
		fmt.Printf("Error counting payment_slips: %v\n", err)
	} else {
		fmt.Printf("Payment slips in PostgreSQL: %d\n", pgSlipsCount)
	}

	var pgPassportsCount int
	err = GetPostgresDB().QueryRowContext(ctx, "SELECT COUNT(*) FROM passenger_passports WHERE deleted_at IS NULL").Scan(&pgPassportsCount)
	if err != nil {
		fmt.Printf("Error counting passenger_passports: %v\n", err)
	} else {
		fmt.Printf("Passenger passports in PostgreSQL: %d\n", pgPassportsCount)
	}

	return nil
}
