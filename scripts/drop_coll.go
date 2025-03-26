package main

import (
	"context"
	"fmt"
	"football_tgbot/db"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set in the .env file")
	}

	client, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("football")
	err = db.Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Database %s dropped successfully\n", "football")
}
