package cache

import (
	"context"
	"time"
)

type UserCache interface {
	Get(ctx context.Context, key string) ([]byte, error)                        // Получить данные из кэша
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error // Записать данные в кэш
	Delete(ctx context.Context, key string) error                               // Удалить данные из кэша
}
