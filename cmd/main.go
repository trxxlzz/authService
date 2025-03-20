package main

import (
	"authService/internal/client/db/pg"
	"authService/internal/config"
	"authService/internal/infra/postgres"
	userRepoPkg "authService/internal/repository/postgres/user"
	//mongoRepoPkg "authService/internal/repository/mongo/user"
	userServPkg "authService/internal/service/user"
	"context"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"

	userAPIPkg "authService/internal/transport/user"
	pb "authService/pkg/protos/gen/go"
)

func main() {
	//Создаем контекст для постгрес
	ctx := context.Background()

	//// Создаём контекст с таймаутом
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	// Загружаем конфиг
	cfg, err := config.LoadConfig("config/.env", "postgres")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	//// Подключаемся к PostreSQL
	dbpool, err := postgres.NewDBConnection(ctx, cfg.DSN("postgres"))
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer dbpool.Close()

	//// Подключаемся к MongoDB
	//client, err := mongoInfra.NewMongoConnection(ctx, cfg.DSN("mongo"))
	//if err != nil {
	//	log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	//}
	//defer func() {
	//	if err := client.Disconnect(ctx); err != nil {
	//		log.Fatalf("Ошибка при отключении от MongoDB: %v", err)
	//	}
	//}()

	log.Println("Successfully connected to database")

	//Инжектим для PostgreSQL
	db := pg.NewDB(dbpool)

	userRepo := userRepoPkg.NewRepository(db)
	userServ := userServPkg.NewService(userRepo)

	////Инжектим для mongoDB
	//userRepo := mongoRepoPkg.NewMongoRepository(client, cfg.MongoDB, "users")
	//userServ := userServPkg.NewService(userRepo)
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
