package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	apiClient "github.com/vsespontanno/tgbot_fschedule/internal/client"
	"github.com/vsespontanno/tgbot_fschedule/internal/db"
	mongorepo "github.com/vsespontanno/tgbot_fschedule/internal/repository/mongodb"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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
	footallClient := apiClient.NewFootballAPIClient(http.DefaultClient, apiKey)
	matchesService := service.NewMatchesService(matchesStore, footallClient)

	calculator := service.NewCalculatorAdapter(teamsStore, standingsStore, matchesStore)
	// Получаем исторические матчи (с 2025-01-01 по 2025-05-06)
	logrus.Info("Fetching historical matches...")
	historicalMatches, err := getHistoricalMatches(matchesService)
	if err != nil {
		logrus.Warnf("Warning: Error fetching historical matches: %v", err)
	} else {
		logrus.Infof("Successfully fetched %d historical matches", len(historicalMatches))
	}

	matches, err := footallClient.FetchMatches(ctx, from, to)
	if err != nil {
		log.Fatal(err)
	}
	err = matchesService.HandleSaveMatches(matches, from, to)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully saved %d matches\n", len(matches))
	for _, match := range matches {
		rating, err := matchesService.CalculateRatingOfMatch(ctx, match, calculator)
		if err != nil {
			logrus.Warnf("Error calculating rating for match %v vs %v; error: %v; skipping", match.HomeTeam.Name, match.AwayTeam.Name, err)
			continue

		}
		match.Rating = rating
		err = matchesService.HandleSaveMatchRating(ctx, match, rating)
		if err != nil {
			logrus.Errorf("Error updating match rating for match %v; error: %v", match, err)
		}
	}

	mongoClient.Disconnect(context.TODO())
}

// Функция для получения исторических матчей с 2025-01-01 по 2025-05-06
func getHistoricalMatches(matchesService *service.MatchesService) ([]types.Match, error) {
	startDate := "2023-06-06"
	endDate := "2025-07-02"

	var allMatches []types.Match
	ctx := context.Background()

	// Разбиваем период на 10-дневные интервалы
	currentDate := startDate
	for currentDate <= endDate {
		// Вычисляем дату окончания для текущего интервала (максимум 10 дней)
		intervalEndDate := addDays(currentDate, 9)
		if intervalEndDate > endDate {
			intervalEndDate = endDate
		}

		logrus.Infof("Fetching matches from %s to %s...", currentDate, intervalEndDate)

		matches, err := matchesService.HandleReqMatches(ctx, currentDate, intervalEndDate)
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
