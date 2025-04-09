package user

import (
	"authService/internal/auth/authRepository"
	"authService/internal/client/db"
	"authService/internal/model"
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
)

type repo struct {
	DB db.Client
}

func NewUserRepo(db db.Client) authRepository.AuthRepository {
	return &repo{DB: db}
}

func (r *repo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := psql.
		Select("*").
		From("users").
		Where(squirrel.Eq{"name": username})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetUserByUsername",
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

	return ToUserFromRepo(&user), nil
}
