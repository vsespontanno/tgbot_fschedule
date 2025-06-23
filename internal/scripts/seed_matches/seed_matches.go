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
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

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

	from := "2025-05-07"
	to := "2025-05-14"
	logrus.Infof("Fetching matches from %s to %s…", from, to)

	mongoURI := os.Getenv("MONGODB_URI")

	// Подключение к MongoDB
	mongoClient, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	matchesStore := mongorepo.NewMongoDBMatchesStore(mongoClient, "football")
	standingsStore := mongorepo.NewMongoDBStandingsStore(mongoClient, "football")
	teamsStore := mongorepo.NewMongoDBTeamsStore(mongoClient, "football")
	standingsService := service.NewStandingService(standingsStore)
	teamsService := service.NewTeamsService(teamsStore)
	matchesService := service.NewMatchesService(matchesStore)

	httpclient := &http.Client{}

	// Получаем исторические матчи (с 2025-01-01 по 2025-05-06)
	logrus.Info("Fetching historical matches...")
	historicalMatches, err := getHistoricalMatches(apiKey, httpclient, mongoClient, matchesService)
	if err != nil {
		logrus.Warnf("Warning: Error fetching historical matches: %v", err)
	} else {
		logrus.Infof("Successfully fetched %d historical matches", len(historicalMatches))
	}

	matches, err := matchesService.HandleReqMatches(httpclient, apiKey, from, to)
	if err != nil {
		log.Fatal(err)
	}
	err = matchesService.HandleSaveMatches(matches, from, to)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully saved %d matches\n", len(matches))
	for _, match := range matches {
		rating, err := giveRatingToMatch(ctx, match, teamsService, standingsService, matchesStore)
		if err != nil {
			logrus.Warnf("Error calculating rating for match %v vs %v; error: %v; skipping", match.HomeTeam.Name, match.AwayTeam.Name, err)
			continue

		}
		match.Rating = rating
		err = updateMatchRatingInMongoDB(mongoClient, match, rating)
		if err != nil {
			logrus.Errorf("Error updating match rating for match %v; error: %v", match, err)
		}
	}

	mongoClient.Disconnect(context.TODO())
}

func giveRatingToMatch(ctx context.Context, match types.Match, teamsService *service.TeamsService, standingsService *service.StandingsService, matchesStore mongorepo.MatchesStore) (float64, error) {

	rating, err := rating.CalculateRatingOfMatch(ctx, match, teamsService, standingsService, matchesStore)
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

// Функция для получения исторических матчей с 2025-01-01 по 2025-05-06
func getHistoricalMatches(apiKey string, httpclient *http.Client, mongoClient *mongo.Client, matchesService *service.MatchesService) ([]types.Match, error) {
	startDate := "2025-01-01"
	endDate := "2025-05-06"

	var allMatches []types.Match

	// Разбиваем период на 10-дневные интервалы
	currentDate := startDate
	for currentDate <= endDate {
		// Вычисляем дату окончания для текущего интервала (максимум 10 дней)
		intervalEndDate := addDays(currentDate, 9)
		if intervalEndDate > endDate {
			intervalEndDate = endDate
		}

		logrus.Infof("Fetching matches from %s to %s...", currentDate, intervalEndDate)

		matches, err := fetchMatchesForPeriod(apiKey, currentDate, intervalEndDate, httpclient)
		if err != nil {
			return nil, fmt.Errorf("error fetching matches for period %s to %s: %w", currentDate, intervalEndDate, err)
		}

		allMatches = append(allMatches, matches...)
		logrus.Infof("Fetched %d matches for period %s to %s", len(matches), currentDate, intervalEndDate)

		// Если это не последний интервал, ждем 10 секунд
		if intervalEndDate < endDate {
			logrus.Infof("Waiting 10 seconds before next request...")
			time.Sleep(10 * time.Second)
		}

		// Переходим к следующему интервалу
		currentDate = addDays(intervalEndDate, 1)
	}

	// Сохраняем все матчи в MongoDB
	if len(allMatches) > 0 {
		err := matchesService.HandleSaveMatches(allMatches, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("error saving historical matches: %w", err)
		}
		logrus.Infof("Successfully saved %d historical matches", len(allMatches))
	}

	return allMatches, nil
}

// Функция для получения матчей за конкретный период
func fetchMatchesForPeriod(apiKey string, startDate string, endDate string, httpclient *http.Client) ([]types.Match, error) {
	url := fmt.Sprintf("https://api.football-data.org/v4/matches?dateFrom=%s&dateTo=%s", startDate, endDate)

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

	// Применяем маппинг названий команд и лиг
	for i := range MatchesResponse.Matches {
		// Маппинг для домашних команд
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

		// Маппинг для гостевых команд
		switch MatchesResponse.Matches[i].AwayTeam.Name {
		case "Wolverhampton Wanderers FC":
			MatchesResponse.Matches[i].AwayTeam.Name = "Wolverhampton FC"
		case "Borussia Mönchengladbach":
			MatchesResponse.Matches[i].AwayTeam.Name = "Borussia Gladbach"
		case "FC Internazionale Milano":
			MatchesResponse.Matches[i].AwayTeam.Name = "Inter"
		case "Club Atlético de Madrid":
			MatchesResponse.Matches[i].AwayTeam.Name = "Atletico Madrid"
		case "RCD Espanyol de Barcelona":
			MatchesResponse.Matches[i].AwayTeam.Name = "Espanyol"
		case "Rayo Vallecano de Madrid":
			MatchesResponse.Matches[i].AwayTeam.Name = "Rayo Vallecano"
		case "Real Betis Balompié":
			MatchesResponse.Matches[i].AwayTeam.Name = "Real Betis"
		case "Real Sociedad de Fútbol":
			MatchesResponse.Matches[i].AwayTeam.Name = "Real Sociedad"
		}

		// Маппинг для названий лиг
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

// Вспомогательная функция для добавления дней к дате
func addDays(dateStr string, days int) string {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		logrus.Infof("Error parsing date %s: %v", dateStr, err)
		return dateStr
	}

	newDate := date.AddDate(0, 0, days)
	return newDate.Format("2006-01-02")
}
