package user

import (
	"microservices/authService/internal/service"
	access "microservices/authService/pkg/jwt/gen/go/access"
	"microservices/authService/pkg/jwt/gen/go/auth"
	user "microservices/authService/pkg/protos/gen/go"
)

type userAPI struct {
	user.UnimplementedUserApiServer
	userService service.UserService
}

type authAPI struct {
	auth.UnimplementedAuthApiServer
	authService service.AuthService
}

type accessApi struct {
	access.UnimplementedAccessApiServer
	accessService service.AccessService
}

func NewUserImplementation(userService service.UserService) *userAPI {
	return &userAPI{userService: userService}
}

func NewAuthImplementation(authService service.AuthService) *authAPI {
	return &authAPI{authService: authService}
}

func NewAccessImplementation(accessService service.AccessService) *accessApi {
	return &accessApi{accessService: accessService}
}
