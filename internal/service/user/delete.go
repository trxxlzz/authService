package user

import (
	"context"
	"fmt"
)

func (s *serv) DeleteUser(ctx context.Context, id int64) error {
	err := s.userRepository.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	// Удаляем кэш
	cacheKey := fmt.Sprintf("user:%d", id)
	_ = s.cache.Delete(ctx, cacheKey)

	return nil
}
