package user

import (
	"authService/internal/model"
	"context"
)

func (s *serv) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.userRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
