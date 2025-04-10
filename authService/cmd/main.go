package main

import (
	"context"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	cache2 "microservices/authService/internal/cache"
	"microservices/authService/internal/client/db/pg"
	"microservices/authService/internal/config"
	"microservices/authService/internal/infra/postgres"
	redisInfra "microservices/authService/internal/infra/redis"
	"microservices/authService/pkg/jwt/gen/go/access"

	"microservices/authService/internal/metric"
	authRepoPkg "microservices/authService/internal/repository/postgres/user"
	userRepoPkg "microservices/authService/internal/repository/postgres/user"
	accessServPkg "microservices/authService/internal/service/user"
	authServPkg "microservices/authService/internal/service/user"
	userServPkg "microservices/authService/internal/service/user"
	accessAPIPkg "microservices/authService/internal/transport/user"
	authAPIPkg "microservices/authService/internal/transport/user"
	userAPIPkg "microservices/authService/internal/transport/user"
	"microservices/authService/pkg/jwt/gen/go/auth"
	user "microservices/authService/pkg/protos/gen/go"
	"net"
)

func main() {
	//Создаем контекст для постгрес
	ctx := context.Background()

	err := metric.Init(ctx)
	if err != nil {
		log.Fatalf("failed to init metrics: %v", err)
	}

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

	// Подключаемся к Redis
	redisClient, err := redisInfra.NewRedisClient(ctx, cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, 0)
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	defer redisClient.Close()

	log.Println("Successfully connected to Redis")

	cache := cache2.NewRedisCache(redisClient)
	db := pg.NewDB(dbpool)

	userRepo := userRepoPkg.NewUserRepository(db)
	userServ := userServPkg.NewUserService(userRepo, cache)

	authRepo := authRepoPkg.NewAuthRepository(db)
	authServ := authServPkg.NewAuthService(authRepo)

	accessServ := accessServPkg.NewAccessService()

	////Инжектим для mongoDB
	//userRepo := mongoRepoPkg.NewMongoRepository(client, cfg.MongoDB, "users")
	//userServ := userServPkg.NewService(userRepo)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051", err)
	}

	//creds, err := credentials.NewServerTLSFromFile("service.pem", "service.key")
	//if err != nil {
	//	log.Fatalf("failed to load TLS keys: %v", err)
	//}

	//основной сервер
	grpcServer := grpc.NewServer(
	//grpc.UnaryInterceptor(
	//interceptor.MetricsInterceptor,
	//),
	) /*grpc.Creds(creds) - тут tls соединение*/

	user.RegisterUserApiServer(grpcServer, userAPIPkg.NewUserImplementation(userServ))
	auth.RegisterAuthApiServer(grpcServer, authAPIPkg.NewAuthImplementation(authServ))
	access.RegisterAccessApiServer(grpcServer, accessAPIPkg.NewAccessImplementation(accessServ))

	lis, err = net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052", err)
	}

	log.Printf("gRPC server is running on port 50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}

	//go func() {
	//	err = runPrometheus()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}()

	//go func() {
	//	log.Printf("gRPC Auth server is running on port 50053")
	//	if err := grpcServerAuth.Serve(lisAuth); err != nil {
	//		log.Fatalf("Failed to serve: %v", err)
	//	}
	//}()

	//тестим access
	//grpcServerAccess := grpc.NewServer()
	//access.RegisterAccessServer(grpcServerAccess, &serverAccess{})

	//lisAccess, err := net.Listen("tcp", ":50054")
	//if err != nil {
	//	log.Fatalf("Failed to listen on port 50054: %v", err)
	//}
	//
	//go func() {
	//	log.Printf("gRPC Access server is running on port 50054")
	//	if err := grpcServerAccess.Serve(lisAccess); err != nil {
	//		log.Fatalf("Failed to serve Access gRPC: %v", err)
	//	}
	//}()
	//
	//select {}

}
