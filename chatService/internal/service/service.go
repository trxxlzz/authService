package service

import (
	"context"
	"microservices/chatService/internal/model"
)

type ChatService interface {
	CreateChat(ctx context.Context, chat *model.Chat) (int64, error)
	DeleteChat(ctx context.Context, chatID int64) error
	SendMessage(ctx context.Context, req *model.Message) error
}
