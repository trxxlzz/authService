package user

import (
	"context"

	"authService/internal/converter"
	pb "authService/pkg/protos/gen/go"
)

func (i *Implementation) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Вызываем метод GetUser у репозитория
	user, err := i.userService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromService(user), nil
}
