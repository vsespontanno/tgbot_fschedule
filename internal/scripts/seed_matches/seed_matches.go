package main

import (
	"context"
	"encoding/json"
	"fmt"
	"football_tgbot/internal/db"
	"football_tgbot/internal/rating"
	mongorepo "football_tgbot/internal/repository/mongodb"
	"football_tgbot/internal/service"
	"football_tgbot/internal/types"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
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

	today := "2025-05-07"
	to := "2025-05-14"
	log.Printf("Fetching matches from %s to %s…", today, to)

	mongoURI := os.Getenv("MONGODB_URI")

	// Подключение к MongoDB
	mongoClient, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	standingsStore := mongorepo.NewMongoDBStandingsStore(mongoClient, "football")
	teamsStore := mongorepo.NewMongoDBTeamsStore(mongoClient, "football")
	standingsService := service.NewStandingService(standingsStore)
	teamsService := service.NewTeamsService(teamsStore)

	httpclient := &http.Client{}

	matches, err := getMatchesSchedule(apiKey, today, to, httpclient, mongoClient)
	if err != nil {
		log.Fatal(err)
	}
	err = saveMatchesToMongoDB(mongoClient, matches, today)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully saved %d matches\n", len(matches))
	for _, match := range matches {
		fmt.Println(match.HomeTeam.Name, match.AwayTeam.Name)
		rating, err := giveRatingToMatch(ctx, match, teamsService, standingsService)
		if err != nil {
			log.Printf("Error calculating rating for match %v vs %v; error: %v; skipping\n", match.HomeTeam.Name, match.AwayTeam.Name, err)
			continue

		}
		match.Rating = rating
		err = updateMatchRatingInMongoDB(mongoClient, match, rating)
		if err != nil {
			log.Printf("Error updating match rating for match %v; error: %v \n", match, err)
		}
	}

	mongoClient.Disconnect(context.TODO())
}

func getMatchesSchedule(apiKey string, today string, tomorrow string, httpclient *http.Client, mongoClient *mongo.Client) ([]types.Match, error) {

	url := fmt.Sprintf("https://api.football-data.org/v4/matches?dateFrom=%s&dateTo=%s", today, tomorrow)
	// req, err := http.NewRequest("GET", today)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Auth-Token", apiKey)

	resp, err := httpclient.Do(req)
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

	leaguesSet := make(map[string]struct{})
	for _, match := range MatchesResponse.Matches {
		leaguesSet[match.Competition.Name] = struct{}{}
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
		case "Primera Division":
			MatchesResponse.Matches[i].Competition.Name = "LaLiga"
		case "Primeira Liga":
			MatchesResponse.Matches[i].Competition.Name = "Primeira"
		case "Premier League":
			MatchesResponse.Matches[i].Competition.Name = "EPL"
		case "Serie A":
			MatchesResponse.Matches[i].Competition.Name = "SerieA"
		case "Bundesliga":
			MatchesResponse.Matches[i].Competition.Name = "Bundesliga"
		case "Ligue 1":
			MatchesResponse.Matches[i].Competition.Name = "Ligue1"
		case "Eredivisie":
			MatchesResponse.Matches[i].Competition.Name = "Eredivisie"
		}
	}
	leaguesSet2 := make(map[string]struct{})
	for _, match := range MatchesResponse.Matches {
		leaguesSet2[match.Competition.Name] = struct{}{}
	}

	// Фильтруем матчи только нужных лиг
	var filteredMatches []types.Match
	allowedLeagues := map[string]bool{
		"LaLiga":     true,
		"EPL":        true,
		"Bundesliga": true,
		"SerieA":     true,
		"Ligue1":     true,
		"UCL":        true,
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

func giveRatingToMatch(ctx context.Context, match types.Match, teamsService *service.TeamsService, standingsService *service.StandingsService) (float64, error) {

	rating, err := rating.CalculateRatingOfMatch(ctx, match, teamsService, standingsService)
	if err != nil {
		return 0, fmt.Errorf("error calculating rating for match %d: %w", match.ID, err)
	}
	return rating, nil

}

func updateMatchRatingInMongoDB(client *mongo.Client, match types.Match, rating float64) error {
	collection := client.Database("football").Collection("matches")

	filter := bson.M{"id": match.ID}
	update := bson.M{"$set": bson.M{"rating": rating}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("error updating match rating for ID %d: %w", match.ID, err)
	}

	return nil
}
