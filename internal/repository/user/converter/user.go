package convert

import (
	"authService/internal/model"
	"authService/internal/repository/user/models"
)

// Конвертируем строку в число чтобы запихнуть в бд
func ConvertUserRoleToInt(role model.UserRole) int {
	switch role {
	case model.UserRoleUser:
		return 1
	case model.UserRoleAdmin:
		return 2
	default:
		return 0
	}
}

// Конвертация модели User в gRPC-ответ - конвертируем из models в model
func ToUserFromRepo(user *models.User) *model.User {
	if user == nil {
		return nil
	}

	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      model.UserRole(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
