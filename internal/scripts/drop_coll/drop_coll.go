package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/vsespontanno/tgbot_fschedule/internal/db"

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
	fmt.Println("---stage 1")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())
	fmt.Println("---stage 2")
	db := client.Database("football")
	fmt.Println("---stage 3")
	err = db.Collection("matches").Drop(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("---stage 4")

	fmt.Println("---stage 5")

	// standings := [6]string{"Bundesliga_standings", "PremierLeague_standings", "LaLiga_standings", "SerieA_standings", "Ligue1_standings", "ChampionsLeague_standings"}
	// for _, s := range standings {
	// 	err = db.Collection(s).Drop(context.TODO())
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
}
