package repository

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"authService/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUser(ctx context.Context, userID int64) (*model.User, error)
	UpdateUser(ctx context.Context, id int64, name string, email string) (*empty.Empty, error)
	DeleteUser(ctx context.Context, id int64) (*empty.Empty, error)
}
