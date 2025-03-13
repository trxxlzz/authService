package user

import (
	"context"

	"authService/internal/model"
)

func (s *serv) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	id, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}
