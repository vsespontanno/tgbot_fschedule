package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"github.com/vsespontanno/tgbot_fschedule/internal/cache"
	"github.com/vsespontanno/tgbot_fschedule/internal/infrastructure/api"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"
	"github.com/vsespontanno/tgbot_fschedule/internal/tools"
)

// Функция, которая апдейтит матче в фоне, пока работает бот
// Работает раз в 24 часа, можно и реже, но пока так.
func RegisterMatchesJob(s *gocron.Scheduler, service *service.MatchesService, redisClient *cache.RedisClient, apiService api.MatchApiClient) {
	logrus.Info("registering matches")
	// _, err := s.Every(24).Hours().Do(func() {
	_, err := s.Every(1).Minute().Do(func() {

		log.Println("Starting matches update...")
		ctx := context.Background()
		start := time.Now()

		from := "2025-05-27"
		to := "2025-06-03"
		matches, err := apiService.FetchMatches(ctx, from, to)
		if err != nil {
			log.Printf("Failed to fetch matches: %v", err)
			return
		}
		tools.MatchFilter(matches)
		log.Printf("Fetched %d matches", len(matches))

		for _, match := range matches {
			fmt.Println(match.HomeTeam.Name + " vs " + match.AwayTeam.Name)
			err = service.HandleUpsertMatch(ctx, match)
			if err != nil {
				log.Printf("Failed to upsert match: %v", err)
			}
		}
		//Очищаем буфер изображений
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
