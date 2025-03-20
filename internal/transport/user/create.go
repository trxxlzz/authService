package user

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"authService/internal/mapper"
	pb "authService/pkg/protos/gen/go"
)

func (i *Implementation) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	id, err := i.userService.CreateUser(ctx, mapper.ToUserFromAPI(req))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateUserResponse{Id: id}, nil
}
