package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (s *serv) UpdateUser(ctx context.Context, id int64, name string, email string) error {
	err := s.userRepository.UpdateUser(ctx, id, name, email)
	if err != nil {
		return err
	}

	// Удаляем старый кэш
	cacheKey := fmt.Sprintf("user:%d", id)
	_ = s.cache.Delete(ctx, cacheKey)

	// Берём обновлённые данные из БД
	updatedUser, err := s.GetUser(ctx, id)
	if err != nil {
		return err
	}

	// Обновляем кэш новыми данными
	userJSON, _ := json.Marshal(updatedUser)
	_ = s.cache.Set(ctx, cacheKey, userJSON, 10*time.Minute)

	return nil
}
