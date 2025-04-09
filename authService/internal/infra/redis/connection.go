package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient создает и возвращает новый клиент Redis
func NewRedisClient(ctx context.Context, host, port, password string, db int) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", host, port) // Формируем адрес

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Проверяем соединение
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("не удалось подключиться к Redis: %v", err)
	}

	return client, nil
}
