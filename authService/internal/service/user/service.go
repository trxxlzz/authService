package user

import (
	"microservices/authService/internal/cache"
	"microservices/authService/internal/repository"
	def "microservices/authService/internal/service"
)

var _ def.UserService = (*userApi)(nil)

type userApi struct {
	userRepository repository.UserRepository
	cache          cache.UserCache
}

type authApi struct {
	authRepository repository.AuthRepository
}

type accessApi struct {
}

func NewUserService(userRepository repository.UserRepository, cache cache.UserCache) *userApi {
	return &userApi{
		userRepository: userRepository,
		cache:          cache,
	}
}

func NewAuthService(userRepository repository.AuthRepository) *authApi {
	return &authApi{
		authRepository: userRepository,
	}
}

func NewAccessService() *accessApi {
	return &accessApi{}
}
