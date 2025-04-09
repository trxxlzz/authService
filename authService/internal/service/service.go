package service

import (
	"context"

	"microservices/authService/internal/model"
)

//go:generate powershell -Command "Remove-Item -Recurse -Force mocks; New-Item -ItemType Directory -Path mocks"
//go:generate minimock -i UserService -o ./mocks/ -s "_minimock.go"

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUser(ctx context.Context, userID int64) (*model.User, error)
	UpdateUser(ctx context.Context, id int64, name string, email string) error
	DeleteUser(ctx context.Context, id int64) error
}

type AuthService interface {
	Login(ctx context.Context, username string, password string) (string, error)
	GetRefreshToken(ctx context.Context, oldToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type AccessService interface {
	CheckAccess(ctx context.Context, endpointAddress string) error
}
