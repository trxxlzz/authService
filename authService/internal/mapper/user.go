package mapper

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"microservices/authService/internal/model"
	auth "microservices/authService/pkg/jwt/gen/go/auth"
	pb "microservices/authService/pkg/protos/gen/go"
	"time"
)

// Конвертирует models.User в pb.GetUserResponse в API слое
func ToUserFromService(user *model.User) *pb.GetUserResponse {
	if user == nil {
		return nil
	}

	return &pb.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      ConvertUserRole(user.Role),
		CreatedAt: timeToProto(user.CreatedAt),
		UpdatedAt: timeToProto(user.UpdatedAt),
	}
}

func ToLoginResponse(refreshToken string) *auth.LoginResponse {
	return &auth.LoginResponse{
		RefreshToken: refreshToken,
	}
}

// timeToProto конвертирует time.Time в *timestamppb.Timestamp
func timeToProto(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

// convertUserRole конвертирует models.UserRole в pb.UserRole
func ConvertUserRole(role model.UserRole) pb.UserRole {
	switch role {
	case model.UserRoleAdmin:
		return pb.UserRole_USER_ROLE_ADMIN
	case model.UserRoleUser:
		return pb.UserRole_USER_ROLE_USER
	default:
		return pb.UserRole_USER_ROLE_UNSPECIFIED // если значение неизвестно
	}
}

// Конвертируем протобаф в models
func ToUserFromAPI(user *pb.CreateUserRequest) *model.User {
	return &model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     model.UserRole(user.Role), // Нужно преобразовать UserRole из pb в model
	}
}
