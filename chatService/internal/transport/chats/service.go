package chats

import (
	"microservices/chatService/internal/service"
	chat "microservices/chatService/pkg/protos/gen/go"
)

type Implementation struct {
	chat.UnimplementedChatApiServer
	chatService service.ChatService
}

func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{chatService: chatService}
}
