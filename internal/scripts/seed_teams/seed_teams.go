package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/vsespontanno/tgbot_fschedule/internal/db"
	"github.com/vsespontanno/tgbot_fschedule/internal/infrastructure/api"
	mongoRepo "github.com/vsespontanno/tgbot_fschedule/internal/repository/mongodb"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"
	"github.com/vsespontanno/tgbot_fschedule/internal/tools"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set in the .env file")
	}

	apiKey := os.Getenv("FOOTBALL_DATA_API_KEY")
	if apiKey == "" {
		log.Fatal("FOOTBALL_DATA_API_KEY is not set in the .env file")
	}

	httpClient := &http.Client{}

	apiClient := api.NewFootballAPIClient(httpClient, apiKey)

	client, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	teamStore := mongoRepo.NewMongoDBTeamsStore(client, "football")
	teamService := service.NewTeamsService(teamStore)
	// Для каждой лиги получаем команды и сохраняем в MongoDB
	for leagueName, league := range types.Leagues {
		log.Printf("Fetching teams for %s...", leagueName)
		teams, err := apiClient.FetchTeams(ctx, league.Code)
		if err != nil {
			log.Printf("Error fetching teams for %s: %v", leagueName, err)
			continue
		}
		if len(teams) == 0 {
			log.Printf("No teams found for %s\n", leagueName)
			continue
		}
		tools.TeamsFilter(teams, leagueName)

		// SAVING IN THEIR OWN LEAGE
		err = teamService.HandleSaveTeams(ctx, leagueName, teams)
		if err != nil {
			log.Printf("Error saving teams in league for %s: %v\n", leagueName, err)
		}
		// SAVING TO COLLECTION WITH ALL
		if leagueName != "ChampionsLeague" {
			err = teamService.HandleSaveTeams(ctx, "Teams", teams)
			if err != nil {
				log.Printf("Error saving teams in all teams for %s:", err)
			}
		}

		fmt.Printf("Successfully saved %d teams for %s\n", len(teams), leagueName)
	}
}
