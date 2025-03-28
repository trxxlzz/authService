package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"authService/internal/model"
)

func (s *serv) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	cacheKey := fmt.Sprintf("user:%d", userID)

	// 1. Проверяем кэш
	cachedUser, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cachedUser != nil {
		var user model.User
		if err := json.Unmarshal(cachedUser, &user); err == nil {
			return &user, nil //Если в кэше есть — сразу возвращаем
		}
	}

	// 2. Если нет в кэше — берём из БД
	user, err := s.userRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. Сохраняем в кэш
	userJSON, _ := json.Marshal(user)
	_ = s.cache.Set(ctx, cacheKey, userJSON, 10*time.Minute)

	return user, nil
}
