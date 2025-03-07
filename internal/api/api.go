package api

import (
	pb "authService/pkg/protos/gen/go"
	"context"
)

type UserHandler struct {
	pb.UnimplementedUserApiServer
}

func (s *UserHandler) CreateUser(ctx context.Context)
