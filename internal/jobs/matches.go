package jobs

import (
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/vsespontanno/tgbot_fschedule/internal/infrastructure/api"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"
)

func RegisterMatchesJob(s *gocron.Scheduler, service *service.MatchesService, apiService api.MatchApiClient) {
	_, err := s.Every(24).Hours().Do(func() {
		log.Println("Starting matches update...")
		ctx := context.Background()
		start := time.Now()

		from := time.Now()
		to := from.Add(24 * time.Hour * 7)
		matches, err := apiService.FetchMatches(ctx, from.Format("2006-01-02"), to.Format("2006-01-02"))
		if err != nil {
			log.Printf("Failed to fetch matches: %v", err)
			return
		}
		for _, match := range matches {
			err = service.HandleUpsertMatch(ctx, match)
			if err != nil {
				log.Printf("Failed to upsert match: %v", err)
			}
		}
		log.Printf("Updated matches schedule (%d records) in %v", len(matches), time.Since(start))
	})

	if err != nil {
		log.Fatalf("Failed to schedule matches job: %v", err)
	}
}
