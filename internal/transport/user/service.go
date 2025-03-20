package user

import (
	"authService/internal/service"
	pb "authService/pkg/protos/gen/go"
)

type Implementation struct {
	pb.UnimplementedUserApiServer
	userService service.UserService
}

func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{userService: userService}
}
