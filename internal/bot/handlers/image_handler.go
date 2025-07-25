package handlers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vsespontanno/tgbot_fschedule/internal/cache"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
	"github.com/vsespontanno/tgbot_fschedule/internal/utils"

	"github.com/sirupsen/logrus"
)

// GenerateTableImage создает изображение турнирной таблицы и сохраняет его в файл.
// Если изображение есть в кэше Redis, возвращает его.
func GenerateTableImage(data []types.Standing, leagueCode string, filename string, redisClient *cache.RedisClient) error {
	cacheKey := fmt.Sprintf("table_image:%s:%s", leagueCode, filename)
	const cacheTTL = 6 * time.Hour
	ctx := context.Background()

	// ПРоверка кеша
	if cachedImage, err := redisClient.GetBytes(ctx, cacheKey); err == nil {
		logrus.WithField("cache_key", cacheKey).Info("Cache hit for table image")
		return os.WriteFile(filename, cachedImage, 0644)
	} else if errors.Is(err, redis.Nil) {
		logrus.WithField("cache_key", cacheKey).Warn("Cache error: ", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("no standings data provided")
	}

	buf, err := utils.TableImage(data)
	if err != nil {
		return fmt.Errorf("failed to generate table image: %s", err)
	}
	// Кэшируем пикчу
	if err := redisClient.SetBytes(ctx, cacheKey, buf.Bytes(), cacheTTL); err != nil {
		return fmt.Errorf("failed to cache table image: %w", err)
	}

	// Сохраняем изображение
	return os.WriteFile(filename, buf.Bytes(), 0644)
}

// функция для генерации изображения расписания матчей
func GenerateScheduleImage(matches []types.Match, filename string, redisClient *cache.RedisClient) error {
	var cacheKey string
	if filename == "" {
		cacheKey = "top_matches_image"
	} else {
		cacheKey = "all_matches_image" + filename
	}
	const cacheTTL = 6 * time.Hour
	ctx := context.Background()

	buf, err := utils.ScheduleImage(matches)
	if err != nil {
		return fmt.Errorf("failed to generate table image: %s", err)
	}

	// Кэшируем изображение
	if err := redisClient.SetBytes(ctx, cacheKey, buf.Bytes(), cacheTTL); err != nil {
		return fmt.Errorf("failed to cache table image: %w", err)
	}
	return os.WriteFile(filename, buf.Bytes(), 0644)

}
