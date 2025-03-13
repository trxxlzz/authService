package user

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

func (s *serv) DeleteUser(ctx context.Context, id int64) (*empty.Empty, error) {
	_, err := s.userRepository.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
