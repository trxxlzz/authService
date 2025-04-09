package user

import (
	"context"
	"microservices/authService/internal/mapper"
	auth "microservices/authService/pkg/jwt/gen/go/auth"
)

func (a *authAPI) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	//лезем в базу или кеш за данными пользователя
	refreshToken, err := a.authService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return mapper.ToLoginResponse(refreshToken), nil
}

func (a *authAPI) GetRefreshToken(ctx context.Context, req *auth.GetRefreshTokenRequest) (*auth.GetRefreshTokenResponse, error) {
	// Логика транспортного слоя: передаем данные в бизнес-слой для выполнения
	refreshToken, err := a.authService.GetRefreshToken(ctx, req.GetOldRefreshToken())
	if err != nil {
		return nil, err
	}

	// Возвращаем новый refresh token в ответе
	return &auth.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}

func (a *authAPI) GetAccessToken(ctx context.Context, req *auth.GetAccessTokenRequest) (*auth.GetAccessTokenResponse, error) {
	// Логика транспортного слоя: передаем данные в бизнес-слой для выполнения
	accessToken, err := a.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	// Возвращаем новый access token в ответе
	return &auth.GetAccessTokenResponse{AccessToken: accessToken}, nil
}
