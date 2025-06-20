package main

import (
	"context"
	"encoding/json"
	"fmt"
	"football_tgbot/internal/db"
	"football_tgbot/internal/types"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func saveTeamsToMongoDB(collection *mongo.Collection, teams []types.Team) error {
	var documents []interface{}
	for _, team := range teams {
		documents = append(documents, team)
	}

	_, err := collection.InsertMany(context.TODO(), documents)
	return err
}

func getTeamsFromAPI(apiKey, leagueCode string) ([]types.Team, error) {
	url := fmt.Sprintf("https://api.football-data.org/v4/competitions/%s/teams", leagueCode)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Auth-Token", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var teamsResponse types.TeamsResponse
	err = json.Unmarshal(body, &teamsResponse)
	if err != nil {
		return nil, err
	}

	return teamsResponse.Teams, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set in the .env file")
	}

	apiKey := os.Getenv("FOOTBALL_DATA_API_KEY")
	if apiKey == "" {
		log.Fatal("FOOTBALL_DATA_API_KEY is not set in the .env file")
	}

	client, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Для каждой лиги получаем команды и сохраняем в MongoDB
	for leagueName, league := range types.Leagues {
		log.Printf("Fetching teams for %s...", leagueName)
		teams, err := getTeamsFromAPI(apiKey, league.Code)
		if err != nil {
			log.Printf("Error fetching teams for %s: %v", leagueName, err)
			continue
		}
		if len(teams) == 0 {
			log.Printf("No teams found for %s\n", leagueName)
			continue
		}
		for i := range teams {
			switch teams[i].Name {
			case "Sevilla FC":
				teams[i].ShortName = "Sevilla"
			case "Wolverhampton Wanderers FC":
				teams[i].Name = "Wolverhampton FC"
			case "Borussia Mönchengladbach":
				teams[i].Name = "Borussia Gladbach"
			case "FC Internazionale Milano":
				teams[i].Name = "Inter"
			case "Club Atlético de Madrid":
				teams[i].Name = "Atletico Madrid"
			case "RCD Espanyol de Barcelona":
				teams[i].Name = "Espanyol"
			case "Rayo Vallecano de Madrid":
				teams[i].Name = "Rayo Vallecano"
			case "Real Betis Balompié":
				teams[i].Name = "Real Betis"
			case "Real Sociedad de Fútbol":
				teams[i].Name = "Real Sociedad"
			}
		}

		for i := range teams {
			switch teams[i].ShortName {
			case "Sevilla FC":
				teams[i].ShortName = "Sevilla"
			case "Leverkusen":
				teams[i].ShortName = "Bayer"
			case "Dortmund":
				teams[i].ShortName = "Borussia D."
			case "M'gladbach":
				teams[i].ShortName = "Borussia M."
			case "Atleti":
				teams[i].ShortName = "Atletico"
			case "Barça":
				teams[i].ShortName = "Barcelona"
			case "Leganés":
				teams[i].ShortName = "Leganes"
			case "Man United":
				teams[i].ShortName = "Manchester United"
			case "Man City":
				teams[i].ShortName = "Manchester City"

			}
		}

		collection := client.Database("football").Collection(leagueName)
		err = saveTeamsToMongoDB(collection, teams)
		if err != nil {
			log.Printf("Error saving teams for %s: %v\n", leagueName, err)
		}

		fmt.Printf("Successfully saved %d teams for %s\n", len(teams), leagueName)
	}
}
