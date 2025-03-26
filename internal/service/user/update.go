package user

import (
	"context"
)

func (s *serv) UpdateUser(ctx context.Context, id int64, name string, email string) error {
	err := s.userRepository.UpdateUser(ctx, id, name, email)
	if err != nil {
		return err
	}

	return nil
}
