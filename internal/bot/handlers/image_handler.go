package handlers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/vsespontanno/tgbot_fschedule/internal/cache"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
	"github.com/vsespontanno/tgbot_fschedule/internal/utils"

	"github.com/sirupsen/logrus"
)

// GenerateTableImage создает изображение турнирной таблицы и сохраняет его в файл.
// Если изображение есть в кэше Redis, возвращает его.
func GenerateTableImage(data []types.Standing, filename string, redisClient *cache.RedisClient) error {
	cacheKey := "table_image: " + filename
	const cacheTTL = 10 * time.Minute
	ctx := context.Background()

	// ПРоверка кеша
	if cachedImage, err := redisClient.GetBytes(ctx, cacheKey); err == nil {
		logrus.WithField("cache_key", cacheKey).Info("Cache hit for table image")
		return os.WriteFile(filename, cachedImage, 0644)
	} else if err.Error() != fmt.Sprintf("cache miss for key %s", cacheKey) {
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
		logrus.WithField("cache_key", cacheKey).Error("Failed to cache table image: ", err)
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
	const cacheTTL = 10 * time.Minute
	ctx := context.Background()

	// Проверяем кэш
	if cachedImage, err := redisClient.GetBytes(ctx, cacheKey); err == nil {
		logrus.WithField("cache_key", cacheKey).Info("Cache hit for schedule image")
		return os.WriteFile(filename, cachedImage, 0644)

	} else if err.Error() != fmt.Sprintf("cache miss for key %s", cacheKey) {
		logrus.WithField("cache_key", cacheKey).Warn("Cache error: ", err)
	}
	buf, err := utils.ScheduleImage(matches)
	if err != nil {
		return fmt.Errorf("failed to generate table image: %s", err)
	}

	// Кэшируем изображение
	if err := redisClient.SetBytes(ctx, cacheKey, buf.Bytes(), cacheTTL); err != nil {
		logrus.WithField("cache_key", cacheKey).Error("Failed to cache schedule image: ", err)
	}
	return os.WriteFile(filename, buf.Bytes(), 0644)

}
