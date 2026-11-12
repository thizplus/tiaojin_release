//go:build ignore

package main

import (
	"context"
	"fmt"
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

	mongoClient, _ := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	defer mongoClient.Disconnect(ctx)
	db := mongoClient.Database("tiaojin")

	fmt.Println("=== ตรวจสอบความสัมพันธ์ Booking <-> User ===\n")

	// 1. ดึงตัวอย่าง booking 5 รายการ
	fmt.Println("--- ตัวอย่าง Booking 5 รายการ ---")
	bookingCursor, _ := db.Collection("tbl_booking").Find(ctx, bson.M{}, options.Find().SetLimit(5))
	var bookings []bson.M
	bookingCursor.All(ctx, &bookings)

	for i, b := range bookings {
		invoiceID := ""
		if inv, ok := b["invoiceId"].(string); ok {
			invoiceID = inv
		}

		profileID := ""
		if pid, ok := b["profileId"].(primitive.ObjectID); ok {
			profileID = pid.Hex()
		} else if pid, ok := b["profileId"].(string); ok {
			profileID = pid
		}

		cusName := ""
		if cn, ok := b["cusName"].(string); ok {
			cusName = cn
		}

		fmt.Printf("%d. Invoice: %s\n", i+1, invoiceID)
		fmt.Printf("   profileId: %s\n", profileID)
		fmt.Printf("   cusName: %s\n", cusName)

		// หา user ที่ตรงกับ profileId
		if profileID != "" {
			oid, err := primitive.ObjectIDFromHex(profileID)
			if err == nil {
				var user bson.M
				err = db.Collection("tbl_users").FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
				if err == nil {
					email := ""
					if e, ok := user["email"].(string); ok {
						email = e
					}
					hasProfile := "ไม่มี"
					if _, ok := user["profile"].(bson.M); ok {
						hasProfile = "มี"
					}
					fmt.Printf("   -> User email: %s, มี profile: %s\n", email, hasProfile)
				} else {
					fmt.Printf("   -> ไม่พบ User ที่ตรงกับ profileId นี้\n")
				}
			}
		}
		fmt.Println()
	}

	// 2. สถิติ
	fmt.Println("--- สถิติ ---")

	// นับจำนวน booking ทั้งหมด
	totalBookings, _ := db.Collection("tbl_booking").CountDocuments(ctx, bson.M{})
	fmt.Printf("Booking ทั้งหมด: %d\n", totalBookings)

	// นับจำนวน user ทั้งหมด
	totalUsers, _ := db.Collection("tbl_users").CountDocuments(ctx, bson.M{})
	fmt.Printf("User ทั้งหมด: %d\n", totalUsers)

	// นับจำนวน user ที่มี profile
	usersWithProfile, _ := db.Collection("tbl_users").CountDocuments(ctx, bson.M{"profile": bson.M{"$ne": nil}})
	fmt.Printf("User ที่มี profile: %d\n", usersWithProfile)

	// ดึง profileId ที่ unique จาก bookings และเช็คว่ามีใน users กี่ตัว
	fmt.Println("\n--- เช็ค profileId ใน Booking ---")
	pipeline := []bson.M{
		{"$group": bson.M{"_id": "$profileId"}},
	}
	profileCursor, _ := db.Collection("tbl_booking").Aggregate(ctx, pipeline)
	var profiles []bson.M
	profileCursor.All(ctx, &profiles)

	matchCount := 0
	notMatchCount := 0
	for _, p := range profiles {
		if pid, ok := p["_id"].(primitive.ObjectID); ok {
			var user bson.M
			err := db.Collection("tbl_users").FindOne(ctx, bson.M{"_id": pid}).Decode(&user)
			if err == nil {
				matchCount++
			} else {
				notMatchCount++
			}
		}
	}
	fmt.Printf("profileId ที่ unique: %d\n", len(profiles))
	fmt.Printf("พบ User ตรงกัน: %d\n", matchCount)
	fmt.Printf("ไม่พบ User: %d\n", notMatchCount)
}
