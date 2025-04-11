package chats

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	pb "microservices/chatService/pkg/protos/gen/go"
)

func (i *Implementation) DeleteChat(ctx context.Context, req *pb.DeleteChatRequest) (*empty.Empty, error) {
	err := i.chatService.DeleteChat(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
