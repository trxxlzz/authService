package user

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"microservices/authService/internal/client/db"
	"microservices/authService/internal/model"
	"microservices/authService/internal/repository"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type authRepo struct {
	DB db.Client
}

func NewAuthRepository(db db.Client) repository.AuthRepository {
	return &authRepo{DB: db}
}

func (r *authRepo) Login(ctx context.Context, username string) (*model.User, error) {
	query := psql.
		Select("name", "role", "password").
		From("users").
		Where(squirrel.Eq{"name": username})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Login",
		QueryRaw: sqlStr,
	}

	var user model.User

	err = r.DB.ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
