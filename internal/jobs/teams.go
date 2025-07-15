package jobs

import (
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/vsespontanno/tgbot_fschedule/internal/infrastructure/api"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
)

func RegisterTeamsJob(s *gocron.Scheduler, service *service.TeamsService, apiService api.TeamsApiClient) {
	ctx := context.Background()
	// _, err := s.Every(6).Hours().Do(
	_, err := s.Every(1).Minutes().Do(func() {
		log.Println("Starting teams update...")
		start := time.Now()

		for leagueName, league := range types.Leagues {

			teams, err := apiService.FetchTeams(ctx, league.Code)
			if err != nil {
				log.Printf("Failed to fetch teams for %s: %v", leagueName, err)
				continue
			}

			if err := service.HandleSaveTeams(context.Background(), league.CollectionName, teams); err != nil {
				log.Printf("Failed to save teams for %s: %v", leagueName, err)
			} else {
				log.Printf("Updated teams for %s (%d records)", leagueName, len(teams))
			}

			if err := service.HandleSaveTeams(context.Background(), "Teams", teams); err != nil {
				log.Printf("Failed to save teams for %s: %v", leagueName, err)
			} else {
				log.Printf("Updated teams for %s (%d records)", leagueName, len(teams))
			}
		}

		log.Printf("teams update completed in %v", time.Since(start))
	})

	if err != nil {
		log.Fatalf("Failed to schedule teams job: %v", err)
	}
}
