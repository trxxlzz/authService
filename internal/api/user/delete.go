package user

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "authService/pkg/protos/gen/go"
)

func (i *Implementation) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*empty.Empty, error) {
	_, err := i.userService.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
