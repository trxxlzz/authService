package converter

import (
	"microservices/chatService/internal/model"
	pb "microservices/chatService/pkg/protos/gen/go"
)

func ConvertFromApi(req *pb.CreateChatRequest) *model.Chat {
	return &model.Chat{
		Usernames: req.Usernames,
	}
}

func ToUserFromApi(req *pb.SendMessageRequest) *model.Message {
	return &model.Message{
		FromUser:  req.From,
		Text:      req.Text,
		Timestamp: req.Timestamp.AsTime(),
	}
}
