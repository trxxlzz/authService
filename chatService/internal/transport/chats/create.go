package chats

import (
	"context"
	"microservices/chatService/internal/converter"
	pb "microservices/chatService/pkg/protos/gen/go"
)

func (i *Implementation) CreateChat(ctx context.Context, req *pb.CreateChatRequest) (*pb.CreateChatResponse, error) {
	id, err := i.chatService.CreateChat(ctx, converter.ConvertFromApi(req))
	if err != nil {
		return nil, err
	}

	return &pb.CreateChatResponse{Id: id}, nil
}
