package main

import (
	userAPIPkg "authService/internal/api/user"
	"authService/internal/config"
	"authService/internal/infra/postgres"
	userRepoPkg "authService/internal/repository/user"
	userServPkg "authService/internal/service/user"
	pb "authService/pkg/protos/gen/go"
	"context"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig("config/.env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := postgres.NewDBConnection(ctx, cfg.DSN())
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	userRepo := userRepoPkg.NewRepository(db)
	userServ := userServPkg.NewService(userRepo)

	grpcServer := grpc.NewServer()

	pb.RegisterUserApiServer(grpcServer, userAPIPkg.NewImplementation(userServ))

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051", err)
	}

	log.Printf("gRPC server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}

}
