package user

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes/empty"
	//"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4/pgxpool"

	"authService/internal/model"
	"authService/internal/repository"
	convert "authService/internal/repository/user/converter"
	"authService/internal/repository/user/models"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type repo struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{DB: db}
}

func (r *repo) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	// Конвертируем строковое значение роли в целое число
	roleInt := convert.ConvertUserRoleToInt(user.Role)

	query := psql.
		Insert("users").
		Columns("name", "email", "password", "role", "created_at").
		Values(user.Name, user.Email, user.Password, roleInt, "NOW()").
		Suffix("RETURNING id")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	var userID int64
	err = r.DB.QueryRow(ctx, sqlStr, args...).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *repo) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	query := psql.
		Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": userID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.DB.QueryRow(ctx, sqlStr, args...)

	var user models.User
	var role int
	err = row.Scan(&user.ID, &user.Name, &user.Email, &role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return convert.ToUserFromRepo(&user), nil
}

func (r *repo) UpdateUser(ctx context.Context, id int64, name string, email string) (*empty.Empty, error) {
	updateQuery := psql.Update("users").
		Set("name", name).
		Set("email", email).
		Where(squirrel.Eq{"id": id})

	sqlStr, args, err := updateQuery.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.DB.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (r *repo) DeleteUser(ctx context.Context, id int64) (*empty.Empty, error) {
	deleteQuery := psql.Delete("users").
		Where(squirrel.Eq{"id": id})

	sqlStr, args, err := deleteQuery.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.DB.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
