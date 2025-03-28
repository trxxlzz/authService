package user

import (
	"authService/internal/client/db"
	"authService/internal/model"
	"authService/internal/repository"
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type repo struct {
	DB db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{DB: db}
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

	return userID, nil
}

func (r *repo) GetUser(ctx context.Context, userID int64) (*model.User, error) {
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

	return nil
}
