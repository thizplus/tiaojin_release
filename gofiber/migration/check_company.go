//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	mongoClient, _ := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	defer mongoClient.Disconnect(ctx)
	db := mongoClient.Database("tiaojin")

	fmt.Println("=== MongoDB tbl_users profile URLs ===\n")

	cursor, _ := db.Collection("tbl_users").Find(ctx, bson.M{"profile": bson.M{"$ne": nil}}, options.Find().SetLimit(5))
	var users []bson.M
	cursor.All(ctx, &users)

	for i, u := range users {
		fmt.Printf("User %d:\n", i+1)
		if profile, ok := u["profile"].(bson.M); ok {
			if logo, ok := profile["logo"].(string); ok && logo != "" {
				fmt.Printf("  logo: %s\n", truncate(logo, 80))
			}
			if atta, ok := profile["atta"].(string); ok && atta != "" {
				fmt.Printf("  atta: %s\n", truncate(atta, 80))
			}
			if dbd, ok := profile["dbd"].(string); ok && dbd != "" {
				fmt.Printf("  dbd:  %s\n", truncate(dbd, 80))
			}
		}
		fmt.Println()
	}
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
