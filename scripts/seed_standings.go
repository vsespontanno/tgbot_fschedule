package main

import (
	"context"
	"encoding/json"
	"fmt"
	"football_tgbot/db"
	"football_tgbot/handlers"
	"football_tgbot/types"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func getLeagueStandings(apiKey, leagueCode string) ([]types.Standing, error) {
	url := fmt.Sprintf("https://api.football-data.org/v4/competitions/%s/standings", leagueCode)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("X-Auth-Token", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var standingsResponse types.StandingsResponse
	err = json.Unmarshal(body, &standingsResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w, body: %s", err, string(body))
	}

	if len(standingsResponse.Standings) > 0 {
		return standingsResponse.Standings[0].Table, nil
	}
	return nil, fmt.Errorf("no standings found for league code: %s", leagueCode)
}

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

	for leagueName, leagueCode := range db.Leagues {
		var standings []types.Standing
		var err error
		for i := 0; i < 3; i++ { // Retry up to 3 times
			standings, err = getLeagueStandings(apiKey, leagueCode)
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

		for i := range standings {
			switch standings[i].Team.Name {
			case "Wolverhampton Wanderers FC":
				standings[i].Team.Name = "Wolverhampton FC"
			case "FC Internazionale Milano":
				standings[i].Team.Name = "Inter"
			case "Club Atlético de Madrid":
				standings[i].Team.Name = "AtLetico Madrid"
			case "RCD Espanyol de Barcelona":
				standings[i].Team.Name = "Espanyol"
			case "Rayo Vallecano de Madrid":
				standings[i].Team.Name = "Rayo Vallecano"
			case "Real Betis Balompié":
				standings[i].Team.Name = "Real Betis"
			case "Real Sociedad de Fútbol":
				standings[i].Team.Name = "Real Sociedad"
			}
		}

		standingsCollection := client.Database("football").Collection(leagueName + "_standings")
		err = handlers.SaveStandingsToMongoDB(standingsCollection, standings)
		if err != nil {
			log.Printf("Error saving standings for %s: %v\n", leagueName, err)
			continue
		}

		fmt.Printf("Successfully saved standings for %s\n", leagueName)
	}
}
