package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"authService/internal/model"
)

func (s *serv) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	id, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	user.ID = id                //Теперь ID обновлён
	user.CreatedAt = time.Now() //Актуализируем дату создания

	// Ключ для кэша
	cacheKey := fmt.Sprintf("user:%d", id)

	// Кладём нового пользователя в кэш
	userJSON, _ := json.Marshal(user)
	err = s.cache.Set(ctx, cacheKey, userJSON, 10*time.Minute)
	if err != nil {
		return id, fmt.Errorf("failed to save user in cache: %w", err)
	}

	return id, nil
}
