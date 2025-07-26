package jobs

import (
	"context"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vsespontanno/tgbot_fschedule/internal/cache"
	"github.com/vsespontanno/tgbot_fschedule/internal/client"
	"github.com/vsespontanno/tgbot_fschedule/internal/service"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"

	"github.com/go-co-op/gocron"
)

// Функция, которая апдейтит турнирные таблицы в фоне, пока работает бот
// Используется gocron для планирования задач
// Каждые 6 часов выполняет обновление турнирных таблиц
// Получает таблицы из API и сохраняет в базу данных
// Очищает кэш Redis для изображений таблиц после обновления
func RegisterStandingsJob(s *gocron.Scheduler, service *service.StandingsService, redisClient *cache.RedisClient, apiService client.StandingsApiClient) {
	logrus.Info("registering standings")
	ctx := context.Background()
	_, err := s.Every(6).Hours().Do(func() {
		log.Println("Starting standings update...")
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
		//Очищаем буфер изобрадений)

		err := redisClient.DeleteByPattern(ctx, "table_image:*")
		if err != nil {
			log.Printf("failed to delete table images: %s", err)
		}

		log.Printf("Standings update completed in %v", time.Since(start))
	})

	if err != nil {
		log.Fatalf("Failed to schedule standings job: %v", err)
	}
}
