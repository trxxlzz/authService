package user

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/olezhek28/platform_common/pkg/sys"
	"github.com/olezhek28/platform_common/pkg/sys/validate"

	"github.com/olezhek28/platform_common/pkg/sys/codes"

	"microservices/authService/internal/mapper"
	user "microservices/authService/pkg/protos/gen/go"
)

func (i *userAPI) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	id, err := i.userService.CreateUser(ctx, mapper.ToUserFromAPI(req))
	if err != nil {
		return nil, err //status.Error(codes.Internal, err.Error())
	}

	return &user.CreateUserResponse{Id: id}, nil
}

func (i *userAPI) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	err := validate.Validate(
		ctx,
		validateID(req.GetId()),
		// otherValidateID(req.GetId()),
	)
	if err != nil {
		return nil, err
	}

	if req.GetId() > 100 {
		return nil, sys.NewCommonError("id must be less than 100", codes.ResourceExhausted)
	}

	// Вызываем метод GetUser у репозитория
	user, err := i.userService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return mapper.ToUserFromService(user), nil
}

func validateID(id int64) validate.Condition {
	return func(ctx context.Context) error {
		if id <= 0 {
			return validate.NewValidationErrors("id must be greater than 0")
		}
		return nil
	}
}

func (i *userAPI) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*empty.Empty, error) {
	err := i.userService.UpdateUser(ctx, req.GetId(), req.GetName().GetValue(), req.GetEmail().GetValue())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (i *userAPI) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*empty.Empty, error) {
	err := i.userService.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
