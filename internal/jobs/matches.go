package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"github.com/vsespontanno/tgbot_fschedule/internal/cache"
	"github.com/vsespontanno/tgbot_fschedule/internal/client"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"
)

// Функция, которая апдейтит матче в фоне, пока работает бот
// Работает раз в 24 часа, можно и реже, но пока так.
func RegisterMatchesJob(s *gocron.Scheduler, service *service.MatchesService, redisClient *cache.RedisClient, apiService client.MatchApiClient, calculator service.Calculator) {
	logrus.Info("registering matches")
	_, err := s.Every(24).Hours().Do(func() {

		log.Println("Starting matches update...")
		ctx := context.Background()
		start := time.Now()

		from := time.Now()
		to := from.Add(24 * time.Hour)
		matches, err := apiService.FetchMatches(ctx, from.Format("2006-01-02"), to.Format("2006-01-02"))
		if err != nil {
			log.Printf("Failed to fetch matches: %v", err)
			return
		}
		log.Printf("Fetched %d matches", len(matches))

		for _, match := range matches {
			fmt.Println(match.HomeTeam.Name + " vs " + match.AwayTeam.Name)
			rating, err := service.CalculateRatingOfMatch(ctx, match, calculator)
			if err != nil {
				logrus.Warnf("Error calculating rating for match %v vs %v; error: %v; skipping", match.HomeTeam.Name, match.AwayTeam.Name, err)
				continue

			}
			match.Rating = rating
			err = service.HandleUpsertMatch(ctx, match)
			if err != nil {
				log.Printf("Failed to upsert match: %v", err)
			}
		}
		//Очищаем буфер изобрадений
		if err := redisClient.DeleteByPattern(ctx, "top_matches_image"); err != nil {
			log.Printf("Failed to delete top matches: %v", err)
		}
		if err := redisClient.DeleteByPattern(ctx, "all_matches*"); err != nil {
			log.Printf("Failed to delete all matches: %v", err)
		}
		log.Printf("Updated matches schedule (%d records) in %v", len(matches), time.Since(start))
	})

	if err != nil {
		log.Fatalf("Failed to schedule matches job: %v", err)
	}
}
