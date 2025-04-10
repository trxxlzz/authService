package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/crypto/bcrypt"
	"microservices/authService/internal/client/db"
	"microservices/authService/internal/model"
	"microservices/authService/internal/repository"
)

var Psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type userRepo struct {
	DB db.Client
}

func NewUserRepository(db db.Client) repository.UserRepository {
	return &userRepo{DB: db}
}

func (r *userRepo) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	//Создаем спан, чтобы зафиксировать время выполнения
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRepo.CreateUser")
	defer span.Finish()

	// Хешируем пароль перед вставкой
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	// Обновляем user.Password на хеш
	user.Password = string(hashedPassword)

	query := Psql.
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

	return userID, nil
}

func (r *userRepo) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	query := Psql.
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
	return resultUser, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, id int64, name string, email string) error {
	updateQuery := Psql.Update("users").
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

	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, id int64) error {
	deleteQuery := Psql.Delete("users").
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

	return nil
}
