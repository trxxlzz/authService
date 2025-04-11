package repository

import (
	"context"

	"microservices/authService/internal/model"
)

//go:generate powershell -Command "Remove-Item -Recurse -Force mocks; New-Item -ItemType Directory -Path mocks"
//go:generate minimock -i UserRepository -o ./mocks/ -s "_minimock.go"

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUser(ctx context.Context, userID int64) (*model.User, error)
	UpdateUser(ctx context.Context, id int64, name string, email string) error
	DeleteUser(ctx context.Context, id int64) error
	UpdateUserRole(ctx context.Context, id int64, role *model.User) error
}

type AuthRepository interface {
	Login(ctx context.Context, username string) (*model.User, error)
}
