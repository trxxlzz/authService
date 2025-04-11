package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"microservices/chatService/internal/client/db/pg"
	"microservices/chatService/internal/interceptor"
	userRepoPkg "microservices/chatService/internal/repository/postgres/chats"
	userServPkg "microservices/chatService/internal/service/chats"
	userAPIPkg "microservices/chatService/internal/transport/chats"
	"net"

	"microservices/chatService/internal/config"
	"microservices/chatService/internal/infra/postgres"
	chat "microservices/chatService/pkg/protos/gen/go"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig("chatService/config/.env")
	if err != nil {
		log.Fatal(err)
	}

	dbpool, err := postgres.NewDBConnection(ctx, cfg.DSN())
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}
	defer dbpool.Close()

	log.Println("Successfully connected to database")

	db := pg.NewDB(dbpool)

	chatRepo := userRepoPkg.NewRepository(db)
	chatServ := userServPkg.NewService(chatRepo)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthInterceptor),
	)

	chat.RegisterChatApiServer(s, userAPIPkg.NewImplementation(chatServ))

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("Server listening on :50053")
}
