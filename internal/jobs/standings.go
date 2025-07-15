package jobs

import (
	"context"
	"log"
	"time"

	"github.com/vsespontanno/tgbot_fschedule/internal/infrastructure/api"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	"github.com/go-co-op/gocron"
)

func RegisterStandingsJob(s *gocron.Scheduler, service *service.StandingsService, apiService api.StandingsApiClient) {
	ctx := context.Background()
	_, err := s.Every(6).Hours().Do(func() {
		log.Println("⏳ Starting standings update...")
		start := time.Now()

		for leagueName, league := range types.Leagues {

			standings, err := apiService.FetchStandings(ctx, league.Code)
			if err != nil {
				log.Printf("Failed to fetch standings for %s: %v", leagueName, err)
				continue
			}

			if err := service.HandleSaveStandings(context.Background(), league.CollectionName, standings); err != nil {
				log.Printf("Failed to save standings for %s: %v", leagueName, err)
			} else {
				log.Printf("Updated standings for %s (%d records)", leagueName, len(standings))
			}
		}

		log.Printf("Standings update completed in %v", time.Since(start))
	})

	if err != nil {
		log.Fatalf("Failed to schedule standings job: %v", err)
	}
}
