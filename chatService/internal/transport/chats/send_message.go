package chats

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"microservices/chatService/internal/converter"
	pb "microservices/chatService/pkg/protos/gen/go"
)

func (i *Implementation) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*empty.Empty, error) {
	err := i.chatService.SendMessage(ctx, converter.ToUserFromApi(req))
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
