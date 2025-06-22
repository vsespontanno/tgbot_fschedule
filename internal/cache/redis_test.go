package cache

import (
	"context"
	"testing"
	"time"
)

func TestNewRedisConn(t *testing.T) {
	redisClient, err := NewRedisClient("redis://localhost:6379")
	if err != nil {
		t.Errorf("NewRedisConn() error = %v", err)
	}
	var res interface{}
	data := "test data"
	exp := time.Minute * 60
	ctx := context.Background()

	err = redisClient.Set(ctx, data, data, exp)
	if err != nil {
		t.Errorf("Set() error = %v", err)
	}

	err = redisClient.Get(ctx, "test", &res)
	if err != nil {
		t.Logf("Get() err = %v", err)
	} else {
		t.Errorf("Get() error = %v", err)
	}

	err = redisClient.Get(ctx, data, &res)
	if err != nil {
		t.Errorf("Get() error = %v", err)
	} else {
		t.Logf("Get() err = %v", err)
	}

	err = redisClient.Delete(ctx, data)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}
	err = redisClient.Get(ctx, data, &res)
	if err == nil {
		t.Error("Expected error for deleted key, but got nil")
	} else {
		t.Logf("Expected cache miss error: %v", err)
	}

	err = redisClient.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestCacheMiss(t *testing.T) {
	redisClient, err := NewRedisClient("redis://localhost:6379")
	if err != nil {
		t.Errorf("NewRedisConn() error = %v", err)
	}
	defer redisClient.Close()

	var res string
	ctx := context.Background()

	// Пытаемся получить несуществующий ключ
	err = redisClient.Get(ctx, "non_existent_key", &res)
	if err == nil {
		t.Error("Expected error for non-existent key, but got nil")
	} else {
		t.Logf("Expected cache miss error: %v", err)
	}
}
