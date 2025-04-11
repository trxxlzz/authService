package chats

import (
	"microservices/chatService/internal/repository"
	def "microservices/chatService/internal/service"
)

var _ def.ChatService = (*serv)(nil)

type serv struct {
	chatRepository repository.ChatRepository
}

func NewService(chatRepository repository.ChatRepository) *serv {
	return &serv{
		chatRepository: chatRepository,
	}
}
