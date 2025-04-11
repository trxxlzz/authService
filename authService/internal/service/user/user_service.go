package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"microservices/authService/internal/model"
)

func (s *userApi) CreateUser(ctx context.Context, user *model.User) (int64, error) {
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

func (s *userApi) GetUser(ctx context.Context, userID int64) (*model.User, error) {
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

func (s *userApi) UpdateUser(ctx context.Context, id int64, name string, email string) error {
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

func (s *userApi) DeleteUser(ctx context.Context, id int64) error {
	err := s.userRepository.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	// Удаляем кэш
	cacheKey := fmt.Sprintf("user:%d", id)
	_ = s.cache.Delete(ctx, cacheKey)

	return nil
}

func (s *userApi) UpdateUserRole(ctx context.Context, id int64, role *model.User) error {
	// Обновляем роль в БД
	if err := s.userRepository.UpdateUserRole(ctx, id, role); err != nil {
		return fmt.Errorf("DB update failed: %w", err)
	}

	// Инвалидируем кэш
	cacheKey := fmt.Sprintf("user:%d", id)
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		log.Printf("Failed to invalidate cache for user %d: %v", id, err)
	}

	return nil
}
