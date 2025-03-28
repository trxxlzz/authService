package user

import (
	"authService/internal/cache"
	"authService/internal/repository"
	def "authService/internal/service"
)

var _ def.UserService = (*serv)(nil)

type serv struct {
	userRepository repository.UserRepository
	cache          cache.UserCache
}

func NewService(userRepository repository.UserRepository, cache cache.UserCache) *serv {
	return &serv{
		userRepository: userRepository,
		cache:          cache,
	}
}
