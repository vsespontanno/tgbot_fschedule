package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vsespontanno/tgbot_fschedule/internal/api"
	"github.com/vsespontanno/tgbot_fschedule/internal/db"
	mongoRepo "github.com/vsespontanno/tgbot_fschedule/internal/repository/mongodb"
	"github.com/vsespontanno/tgbot_fschedule/internal/tools"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	apiKey := os.Getenv("FOOTBALL_DATA_API_KEY")
	if apiKey == "" {
		log.Fatal("FOOTBALL_DATA_API_KEY is not set in the .env file")
	}
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set in the .env file")
	}

	client, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	store := mongoRepo.NewMongoDBStandingsStore(client, "football")
	footallClient := api.NewFootballAPIClient(&http.Client{}, apiKey)
	for leagueName, league := range types.Leagues {
		var standings []types.Standing
		var err error
		for i := 0; i < 3; i++ { // Retry up to 3 times
			if league.Code == "CL" {
				break
			}
			standings, err = footallClient.FetchStandings(context.Background(), league.Code)
			if err == nil {
				break // Success, exit retry loop
			}
			log.Printf("Attempt %d: Error getting standings for %s: %v\n", i+1, leagueName, err)
			time.Sleep(2 * time.Second) // Wait for 2 seconds before retrying
		}
		if err != nil {
			log.Printf("Failed to get standings for %s after multiple retries: %v\n", leagueName, err)
			continue
		}

		if len(standings) == 0 {
			log.Printf("No standings found for %s\n", leagueName)
			continue
		}

		tools.StandingsFilter(standings)

		err = store.SaveStandings(context.Background(), leagueName, standings)
		if err != nil {
			log.Printf("Error saving standings for %s: %v\n", leagueName, err)
			continue
		}
		fmt.Printf("Saved standings for %s\n", leagueName)

	}
}
