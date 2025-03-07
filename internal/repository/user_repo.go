package repository

import (
	"authService/internal/models"
	pb "authService/pkg/protos/gen/go"
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	query := psql.
		Insert("users").
		Columns("name", "email", "password", "role", "created_at").
		Values(user.Name, user.Email, user.Password, user.Role, "NOW()").
		Suffix("RETURNING id")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	var userID int64
	err = r.DB.QueryRowContext(ctx, sqlStr, args...).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *UserRepository) GetUser(ctx context.Context, userID int64) (*pb.GetUserResponse, error) {
	query := psql.
		Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": userID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.DB.QueryRowContext(ctx, sqlStr, args...)

	var user models.User
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Преобразуем время в google.protobuf.Timestamp
	createdAt := timestamppb.New(user.CreatedAt)
	updatedAt := timestamppb.New(user.UpdatedAt)

	var role pb.UserRole
	switch user.Role {
	case models.UserRoleUser:
		role = pb.UserRole_USER_ROLE_USER
	case models.UserRoleAdmin:
		role = pb.UserRole_USER_ROLE_ADMIN
	default:
		role = pb.UserRole_USER_ROLE_UNSPECIFIED
	}

	// Формируем и возвращаем ответ
	return &pb.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      role, // Преобразуем роль
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id int64, name string, email string) (*empty.Empty, error) {
	updateQuery := psql.Update("users").
		Set("name", name).
		Set("email", email).
		Where(squirrel.Eq{"id": id})

	sqlStr, args, err := updateQuery.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.DB.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int64) (*empty.Empty, error) {
	deleteQuery := psql.Delete("users").
		Where(squirrel.Eq{"id": id})

	sqlStr, args, err := deleteQuery.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.DB.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
