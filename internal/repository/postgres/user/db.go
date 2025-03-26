package user

import (
	"authService/internal/repository"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"

	"authService/internal/client/db"
	"authService/internal/model"
	"github.com/Masterminds/squirrel"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type repo struct {
	DB          db.Client
	RedisClient *redis.Client
}

func NewRepository(db db.Client, redisClient *redis.Client) repository.UserRepository {
	return &repo{
		DB:          db,
		RedisClient: redisClient}
}

func (r *repo) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	query := psql.
		Insert("users").
		Columns("name", "email", "password", "role", "created_at").
		Values(user.Name, user.Email, user.Password, user.Role, "NOW()").
		Suffix("RETURNING id")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: sqlStr,
	}

	var userID int64
	err = r.DB.QueryRowContext(ctx, q, args...).Scan(&userID)
	if err != nil {
		return 0, err
	}

	/// Присваиваем ID пользователю
	user.ID = userID
	user.CreatedAt = time.Now()

	// Сериализуем пользователя в JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal user: %w", err)
	}

	// Ключ для кэша в Redis
	cacheKey := fmt.Sprintf("user:%d", userID)

	// Сохраняем в Redis с TTL 10 минут
	err = r.RedisClient.Set(ctx, cacheKey, userJSON, 10*time.Minute).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to save user in Redis: %w", err)
	}

	return userID, nil
}

func (r *repo) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	cacheKey := fmt.Sprintf("user:%d", userID)

	//Проверяем кеш
	cachedUser, err := r.RedisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var user model.User
		if err := json.Unmarshal(cachedUser, &user); err == nil {
			return &user, nil // ✅ Если в кеше есть — возвращаем
		}
	}

	//Если нет в кеше, берем из PostgreSQL
	query := psql.
		Select("*").
		From("users").
		Where(squirrel.Eq{"id": userID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetUser",
		QueryRaw: sqlStr,
	}

	var user User
	err = r.DB.ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	resultUser := ToUserFromRepo(&user)

	//Кладем в кеш (TTL = 10 минут)
	userJSON, _ := json.Marshal(resultUser)
	r.RedisClient.Set(ctx, cacheKey, userJSON, 10*time.Minute)

	return resultUser, nil
}

func (r *repo) UpdateUser(ctx context.Context, id int64, name string, email string) error {
	updateQuery := psql.Update("users").
		Set("name", name).
		Set("email", email).
		Where(squirrel.Eq{"id": id})

	sqlStr, args, err := updateQuery.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.UpdateUser",
		QueryRaw: sqlStr,
	}

	_, err = r.DB.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	// Удаляем старый кеш перед получением свежих данных
	cacheKey := fmt.Sprintf("user:%d", id)
	err = r.RedisClient.Del(ctx, cacheKey).Err()
	if err != nil {
		return err
	}

	// После обновления забираем новые данные пользователя
	updatedUser, err := r.GetUser(ctx, id)
	if err != nil {
		return err
	}

	// Обновляем кеш новыми данными (TTL = 10 минут)
	cacheKey = fmt.Sprintf("user:%d", id)
	userJSON, _ := json.Marshal(updatedUser)
	err = r.RedisClient.Set(ctx, cacheKey, userJSON, 10*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) DeleteUser(ctx context.Context, id int64) error {
	deleteQuery := psql.Delete("users").
		Where(squirrel.Eq{"id": id})

	sqlStr, args, err := deleteQuery.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.DeleteUser",
		QueryRaw: sqlStr,
	}

	_, err = r.DB.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	// Удаляем кеш, если пользователь был удален
	cacheKey := fmt.Sprintf("user:%d", id)
	r.RedisClient.Del(ctx, cacheKey)

	return nil
}
