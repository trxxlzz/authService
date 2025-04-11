package user

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
	"log"
	"microservices/authService/internal/model"
	"microservices/authService/internal/utils"
	"strings"
)

func (a *accessApi) CheckAccess(ctx context.Context, endpointAddress string) error {
	// Извлечение метаданных из контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("metadata is not provided")
	}

	// Извлечение заголовка authorization
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return errors.New("authorization header is not provided")
	}

	accessToken := strings.TrimSpace(strings.TrimPrefix(authHeader[0], authPrefix))

	// Проверка токена
	claims, err := utils.VerifyToken(accessToken, accessTokenSecretKey)
	if err != nil {
		log.Printf("Token verification failed: %s", err.Error())
		return errors.New("access token is invalid")
	}

	// Получаем доступные роли для эндпоинтов
	accessibleMap, err := a.getAccessibleRoles(ctx)
	if err != nil {
		return errors.New("failed to get accessible roles")
	}

	// Проверка доступа
	role, ok := accessibleMap[endpointAddress]
	if !ok {
		return nil // Нет роли для этого эндпоинта, доступ разрешен
	}

	// Сравнение ролей
	if role == claims.Role {
		return nil
	}

	// Отказ в доступе
	return errors.New("access denied")
}

func (a *accessApi) getAccessibleRoles(ctx context.Context) (map[string]model.UserRole, error) {

	accessRoles = make(map[string]model.UserRole)

	accessRoles[model.CreateChat] = model.UserRoleAdmin

	accessRoles[model.DeleteChat] = model.UserRoleAdmin
	// Можно добавить другие эндпоинты и соответствующие роли

	return accessRoles, nil
}

var accessRoles map[string]model.UserRole
