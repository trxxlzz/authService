package mapper

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"microservices/authService/internal/model"
	auth "microservices/authService/pkg/jwt/gen/go/auth"
	pb "microservices/authService/pkg/protos/gen/go"
	"strconv"
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

func ToUserFromKafka(message *kafka.Message) (int64, *model.User, error) {
	var kafkaMsg struct {
		UserID string `json:"user_id"`
		Role   int    `json:"role"`
	}

	if err := json.Unmarshal(message.Value, &kafkaMsg); err != nil {
		return 0, nil, fmt.Errorf("unmarshal error: %w", err)
	}

	// Конвертируем string UserID в int64
	id, err := strconv.ParseInt(kafkaMsg.UserID, 10, 64)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	userRole := model.UserRole(kafkaMsg.Role)

	// Создаем объект User с обновленной ролью
	user := &model.User{
		ID:        id,
		Role:      userRole,
		UpdatedAt: time.Now(), // Устанавливаем текущее время обновления
	}

	return id, user, nil
}
