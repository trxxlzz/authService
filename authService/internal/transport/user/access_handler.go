package user

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
	"microservices/authService/pkg/jwt/gen/go/access"
)

func (a *accessApi) Check(ctx context.Context, req *access.CheckRequest) (*emptypb.Empty, error) {
	// Передаем запрос в бизнес-слой для проверки доступа
	err := a.accessService.CheckAccess(ctx, req.GetEndpointAddress())
	if err != nil {
		return &empty.Empty{}, err
	}

	return &empty.Empty{}, nil
}
