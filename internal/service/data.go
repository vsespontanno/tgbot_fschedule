package service

import (
	"context"
	"football_tgbot/internal/config"
	"football_tgbot/internal/db"
	mongoRepo "football_tgbot/internal/repository/mongodb"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func UpdateData() {
	// Загрузка конфига
	cfg := config.LoadConfig("../../.env")
	mongoURI := cfg.MongoURI
	footballAPI := cfg.FootballDataAPIKey

	logrus.Info("Updating data")

	// Сегодняшняя дата и конечная
	from := time.Now().Format("2006-01-02")
	to := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	// Подключение к MongoDB
	mongoClient, err := db.ConnectToMongoDB(mongoURI)
	if err != nil {
		logrus.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.TODO())
	logrus.Info("Connected to MongoDB")

	//Подключение к http-клиенту
	httpclient := &http.Client{}

	//Сервис матчей
	matchesStore := mongoRepo.NewMongoDBMatchesStore(mongoClient, "football")
	matchesService := NewMatchesService(matchesStore)

	logrus.Infof("Fetching matches from %s to %s…", from, to)
	matches, err := matchesService.HandleReqMatches(httpclient, footballAPI, from, to)
	if err != nil {
		logrus.Fatalf("Error fetching matches: %v", err)
	}

	for i, v := range matches {
		match, err := matchesService.HandleGetMatchByID(context.Background(), v.ID)
		if err != nil {
			logrus.Fatal("Error while checking if match already in db: ", err)
		} else if match.ID != 0 {
			matches = append(matches[:i], matches[i+1:]...)
		}

	}
	err = matchesService.HandleSaveMatches(matches, from, to)
	if err != nil {
		logrus.Fatalf("Error saving matches: %v", err)
	}
	logrus.Infof("Successfully saved %d matches", len(matches))
}
