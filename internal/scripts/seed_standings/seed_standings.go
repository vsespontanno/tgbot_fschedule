package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vsespontanno/tgbot_fschedule/internal/db"
	mongoRepo "github.com/vsespontanno/tgbot_fschedule/internal/repository/mongodb"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	"github.com/joho/godotenv"
)

func getLeagueStandings(apiKey, leagueCode string) ([]types.Standing, error) {
	url := fmt.Sprintf("https://api.football-data.org/v4/competitions/%s/standings?season=2024", leagueCode)
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
	fmt.Println("Stage 1: Loading environment variables")
	apiKey := os.Getenv("FOOTBALL_DATA_API_KEY")
	if apiKey == "" {
		log.Fatal("FOOTBALL_DATA_API_KEY is not set in the .env file")
	}
	fmt.Printf("Stage 2: Loading apiKey - %s\n", apiKey)
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set in the .env file")
	}
	fmt.Printf("Stage 3: Loading mongoURI - %s\n", mongoURI)

	client, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	fmt.Println("Stage 4: Connected to MongoDB")
	defer client.Disconnect(context.TODO())

	store := mongoRepo.NewMongoDBStandingsStore(client, "football")
	fmt.Println("Stage 5: Created MongoDBStandingsStore")

	for leagueName, league := range types.Leagues {
		var standings []types.Standing
		var err error
		for i := 0; i < 3; i++ { // Retry up to 3 times
			if league.Code == "CL" {
				break
			}
			standings, err = getLeagueStandings(apiKey, league.Code)
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
			case "Borussia Mönchengladbach":
				standings[i].Team.Name = "Borussia Gladbach"
			case "FC Internazionale Milano":
				standings[i].Team.Name = "Inter"
			case "Club Atlético de Madrid":
				standings[i].Team.Name = "Atletico Madrid"
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

		for i := range standings {
			switch standings[i].Team.ShortName {
			case "Leverkusen":
				standings[i].Team.ShortName = "Bayer"
			case "Dortmund":
				standings[i].Team.ShortName = "Borussia D."
			case "M'gladbach":
				standings[i].Team.ShortName = "Borussia M."
			case "Atleti":
				standings[i].Team.ShortName = "Atletico"
			case "Barça":
				standings[i].Team.ShortName = "Barcelona"
			case "Leganés":
				standings[i].Team.ShortName = "Leganes"
			case "Man United":
				standings[i].Team.ShortName = "Manchester United"
			case "Man City":
				standings[i].Team.ShortName = "Manchester City"

			}
		}

		err = store.SaveStandings(context.Background(), leagueName, standings)
		if err != nil {
			log.Printf("Error saving standings for %s: %v\n", leagueName, err)
			continue
		}
		fmt.Printf("Saved standings for %s\n", leagueName)

	}
}
