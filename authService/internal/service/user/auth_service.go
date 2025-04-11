package user

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices/authService/internal/model"
	"microservices/authService/internal/utils"
	"time"
)

const (
	authPrefix = "Bearer"

	refreshTokenSecretKey = "2ltaEBV55ZZ+OFmXzgTlc/qz2OJF5doWW4gbs+pkmP8="
	accessTokenSecretKey  = "cKmx+rugR6xWWTcC4ZHNh+6buhF5LgBtiJ6JECYDv2k="

	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 10 * time.Minute
)

func (a *authApi) Login(ctx context.Context, username, password string) (string, error) {
	//лезем в базу или кеш за данными пользователя
	//сверяем хэши пароля
	user, err := a.authRepository.Login(ctx, username)
	if err != nil {
		return "", err
	}

	if !utils.VerifyPassword(user.Password, password) {
		return "", errors.New("invalid credentials")
	}

	refreshToken, err := utils.GenerateToken(model.User{
		Name: user.Name,
		Role: user.Role,
	},
		refreshTokenSecretKey,
		refreshTokenExpiration,
	)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return refreshToken, nil
}

func (a *authApi) GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	claims, err := utils.VerifyToken(oldRefreshToken, refreshTokenSecretKey)
	if err != nil {
		return "", status.Errorf(codes.Aborted, "invalid refresh token")
	}

	refreshToken, err := utils.GenerateToken(model.User{
		Name: claims.Name,
		Role: claims.Role,
	},
		refreshTokenSecretKey,
		refreshTokenExpiration,
	)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (a *authApi) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, refreshTokenSecretKey)
	if err != nil {
		return "", status.Errorf(codes.Aborted, "incalid refresh token")
	}

	accessToken, err := utils.GenerateToken(model.User{
		Name: claims.Name,
		Role: claims.Role,
	},
		accessTokenSecretKey,
		accessTokenExpiration,
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
