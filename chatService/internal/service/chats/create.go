package chats

import (
	"context"

	"microservices/chatService/internal/model"
)

func (s *serv) CreateChat(ctx context.Context, chat *model.Chat) (int64, error) {
	id, err := s.chatRepository.CreateChat(ctx, chat)
	if err != nil {
		return 0, err
	}

	return id, nil
}
