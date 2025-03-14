package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"

	_ "github.com/lib/pq"

	userAPIPkg "authService/internal/api/user"
	"authService/internal/client/db/pg"
	"authService/internal/config"
	"authService/internal/infra/postgres"
	userRepoPkg "authService/internal/repository/user"
	userServPkg "authService/internal/service/user"
	pb "authService/pkg/protos/gen/go"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig("config/.env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbpool, err := postgres.NewDBConnection(ctx, cfg.DSN())
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer dbpool.Close()

	log.Println("Successfully connected to database")

	db := pg.NewDB(dbpool)

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
