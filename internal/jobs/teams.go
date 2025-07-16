package jobs

import (
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"github.com/vsespontanno/tgbot_fschedule/internal/api"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
)

// Функция, которая апдейтит матче в фоне, пока работает бот
// Работает раз в 300 дней, т.к. команды меняются только после конца сезона:
// Команды с нисших дивизионов попадают в дивизионы повыше и наоборот
func RegisterTeamsJob(s *gocron.Scheduler, service *service.TeamsService, apiService api.TeamsApiClient) {
	logrus.Info("registering teams")

	ctx := context.Background()
	// _, err := s.Every(24 * 10 * 30).Hours().
	_, err := s.Every(1).Minute().Do(func() {
		log.Println("Starting teams update...")
		start := time.Now()

		for leagueName, league := range types.Leagues {

			teams, err := apiService.FetchTeams(ctx, league.Code)
			if err != nil {
				log.Printf("Failed to fetch teams for %s: %v", leagueName, err)
				continue
			}
			for _, team := range teams {
				team.League = leagueName

				if err := service.HandleUpsertMatch(context.Background(), league.CollectionName, team); err != nil {
					log.Printf("Failed to save team for %s: %v", leagueName, err)
				}

				if err := service.HandleUpsertMatch(context.Background(), "Teams", team); err != nil {
					log.Printf("Failed to save teams for %s: %v", leagueName, err)
				}
			}
		}

		log.Printf("teams update completed in %v", time.Since(start))
	})

	if err != nil {
		log.Fatalf("Failed to schedule teams job: %v", err)
	}
}
