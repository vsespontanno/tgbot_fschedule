package jobs

import (
	"context"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vsespontanno/tgbot_fschedule/internal/infrastructure/api"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"
	"github.com/vsespontanno/tgbot_fschedule/internal/tools"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	"github.com/go-co-op/gocron"
)

// Функция, которая апдейтит турнирные таблицы в фоне, пока работает бот
// Работает раз в 6 часов, т.к. после каждого матча таблица обновляется
func RegisterStandingsJob(s *gocron.Scheduler, service *service.StandingsService, apiService api.StandingsApiClient) {
	logrus.Info("registering standings")
	ctx := context.Background()
	// _, err := s.Every(6).Hours().Do(func() {
	_, err := s.Every(1).Minute().Do(func() {

		log.Println("Starting standings update...")
		start := time.Now()

		for leagueName, league := range types.Leagues {

			standings, err := apiService.FetchStandings(ctx, league.Code)
			if err != nil {
				log.Printf("Failed to fetch standings for %s: %v", leagueName, err)
				continue
			}
			tools.StandingsFilter(standings)

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
