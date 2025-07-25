package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// RedisClient обертка для клиента Redis.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient создает новый клиент Redis.
func NewRedisClient(redisURL string) (*RedisClient, error) {
	logrus.Info("Connecting to Redis...")
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logrus.Info("Connected to Redis")

	return &RedisClient{client: client}, nil
}

// SetBytes сохраняет байтовый массив в Redis с указанным временем жизни.
func (c *RedisClient) SetBytes(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// GetBytes получает байтовый массив из Redis по ключу.
// Возвращает ошибку, если ключ не найден или произошла другая ошибка.
func (c *RedisClient) GetBytes(ctx context.Context, key string) ([]byte, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("cache miss for key %s", key)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}
	return data, nil
}

// DeleteByPattern удаляет все ключи, соответствующие шаблону.
// Использует SCAN для безопасного удаления ключей в больших базах данных.
// Возвращает ошибку, если не удалось сканировать ключи или удалить их.
func (c *RedisClient) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan keys: %w", err)
	}
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}

// Close закрывает соединение с Redis.
// Возвращает ошибку, если не удалось закрыть соединение.
// Используется для освобождения ресурсов при завершении работы приложения.
func (c *RedisClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
