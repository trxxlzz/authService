package chats

import (
	"context"

	"microservices/chatService/internal/model"
)

func (s *serv) SendMessage(ctx context.Context, req *model.Message) error {
	err := s.chatRepository.SendMessage(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
