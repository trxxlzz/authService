package user

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "authService/pkg/protos/gen/go"
)

func (i *Implementation) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*empty.Empty, error) {
	_, err := i.userService.UpdateUser(ctx, req.GetId(), req.GetName().GetValue(), req.GetEmail().GetValue())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
