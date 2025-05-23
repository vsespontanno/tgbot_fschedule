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
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Загрузка .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получение значений из .env
	apiKey := os.Getenv("FOOTBALL_DATA_API_KEY")
	if apiKey == "" {
		log.Fatal("FOOTBALL_DATA_API_KEY is not set in the .env file")
	}

	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	mongoURI := os.Getenv("MONGODB_URI")

	// Подключение к MongoDB
	client, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		log.Fatal(err)
	}

	httpclient := &http.Client{}
	matches, err := getMatchesSchedule(apiKey, today, tomorrow, httpclient)
	if err != nil {
		log.Fatal(err)
	}
	err = saveMatchesToMongoDB(client, matches, today)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully saved %d matches\n", len(matches))

	client.Disconnect(context.TODO())
}

func getMatchesSchedule(apiKey string, today string, tomorrow string, client *http.Client) ([]types.Match, error) {
	url := fmt.Sprintf("https://api.football-data.org/v4/matches?dateFrom=%s&dateTo=%s", today, tomorrow)
	// req, err := http.NewRequest("GET", today)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Auth-Token", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var MatchesResponse types.MatchesResponse
	err = json.Unmarshal(body, &MatchesResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}
	for i := range MatchesResponse.Matches {
		switch MatchesResponse.Matches[i].HomeTeam.Name {
		case "Wolverhampton Wanderers FC":
			MatchesResponse.Matches[i].HomeTeam.Name = "Wolverhampton FC"
		case "Borussia Mönchengladbach":
			MatchesResponse.Matches[i].HomeTeam.Name = "Borussia Gladbach"
		case "FC Internazionale Milano":
			MatchesResponse.Matches[i].HomeTeam.Name = "Inter"
		case "Club Atlético de Madrid":
			MatchesResponse.Matches[i].HomeTeam.Name = "Atletico Madrid"
		case "RCD Espanyol de Barcelona":
			MatchesResponse.Matches[i].HomeTeam.Name = "Espanyol"
		case "Rayo Vallecano de Madrid":
			MatchesResponse.Matches[i].HomeTeam.Name = "Rayo Vallecano"
		case "Real Betis Balompié":
			MatchesResponse.Matches[i].HomeTeam.Name = "Real Betis"
		case "Real Sociedad de Fútbol":
			MatchesResponse.Matches[i].HomeTeam.Name = "Real Sociedad"
		}
		// Corrected loop for AwayTeam
		switch MatchesResponse.Matches[i].AwayTeam.Name {
		case "Wolverhampton Wanderers FC":
			MatchesResponse.Matches[i].AwayTeam.Name = "Wolverhampton FC"
		case "Borussia Mönchengladbach":
			MatchesResponse.Matches[i].AwayTeam.Name = "Borussia Gladbach"
		case "FC Internazionale Milano":
			MatchesResponse.Matches[i].AwayTeam.Name = "Inter"
		case "Club Atlético de Madrid":
			MatchesResponse.Matches[i].AwayTeam.Name = "AtLetico Madrid"
		case "RCD Espanyol de Barcelona":
			MatchesResponse.Matches[i].AwayTeam.Name = "Espanyol"
		case "Rayo Vallecano de Madrid":
			MatchesResponse.Matches[i].AwayTeam.Name = "Rayo Vallecano"
		case "Real Betis Balompié":
			MatchesResponse.Matches[i].AwayTeam.Name = "Real Betis"
		case "Real Sociedad de Fútbol":
			MatchesResponse.Matches[i].AwayTeam.Name = "Real Sociedad"
		}

		switch MatchesResponse.Matches[i].Competition.Name {
		case "UEFA Champions League":
			MatchesResponse.Matches[i].Competition.Name = "UCL"
		case "UEFA Europa League":
			MatchesResponse.Matches[i].Competition.Name = "UEL"
		case "Primera División":
			MatchesResponse.Matches[i].Competition.Name = "La Liga"
		case "Primeira Liga":
			MatchesResponse.Matches[i].Competition.Name = "Primeira"
		case "Premier League":
			MatchesResponse.Matches[i].Competition.Name = "EPL"
		case "Serie A":
			MatchesResponse.Matches[i].Competition.Name = "Serie A"
		case "Bundesliga":
			MatchesResponse.Matches[i].Competition.Name = "Bundesliga"
		case "Eredivisie":
			MatchesResponse.Matches[i].Competition.Name = "Eredivisie"
		}

	}

	// Фильтруем матчи только нужных лиг
	var filteredMatches []types.Match
	allowedLeagues := map[string]bool{
		"La Liga":    true,
		"EPL":        true,
		"Primeira":   true,
		"Eredivisie": true,
		"Bundesliga": true,
		"Serie A":    true,
		"UCL":        true,
		"UEL":        true,
	}

	for _, match := range MatchesResponse.Matches {
		if allowedLeagues[match.Competition.Name] {
			filteredMatches = append(filteredMatches, match)
		}
	}

	return filteredMatches, nil
}

func saveMatchesToMongoDB(client *mongo.Client, matches []types.Match, today string) error {
	if len(matches) == 0 {
		log.Printf("No matches found for today: %s\n", today)
		return nil
	}

	collection := client.Database("football").Collection("matches")

	var documents []interface{}
	for _, match := range matches {
		documents = append(documents, match)
	}

	_, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		return fmt.Errorf("error inserting matches: %v", err)
	}

	return nil
}
