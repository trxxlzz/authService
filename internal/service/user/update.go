package user

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *serv) UpdateUser(ctx context.Context, id int64, name string, email string) (*empty.Empty, error) {
	_, err := s.userRepository.UpdateUser(ctx, id, name, email)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
