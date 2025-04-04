package main

import (
	cache2 "authService/internal/cache"
	"authService/internal/client/db/pg"
	"authService/internal/config"
	"authService/internal/infra/postgres"
	"authService/internal/interceptor"
	"authService/internal/metric"
	"authService/internal/model"
	userRepoPkg "authService/internal/repository/postgres/user"
	"authService/internal/utils"
	"authService/pkg/jwt/gen/go/access"
	"authService/pkg/jwt/gen/go/auth"
	"errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"strings"
	"time"

	//mongoRepoPkg "authService/internal/repository/mongo/user"
	redisInfra "authService/internal/infra/redis"
	userServPkg "authService/internal/service/user"
	userAPIPkg "authService/internal/transport/user"
	pb "authService/pkg/protos/gen/go"
	"context"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	authPrefix = "Bearer"

	refreshTokenSecretKey = "2ltaEBV55ZZ+OFmXzgTlc/qz2OJF5doWW4gbs+pkmP8="
	accessTokenSecretKey  = "cKmx+rugR6xWWTcC4ZHNh+6buhF5LgBtiJ6JECYDv2k="

	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 2 * time.Minute
)

type serverAccess struct {
	access.UnimplementedAccessServer
}

type serverAuth struct {
	auth.UnimplementedAuthApiServer
}

func (s *serverAuth) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	//лезем в базу или кеш за данными пользователя
	//сверяем хэши пароля

	refreshToken, err := utils.GenerateToken(model.User{
		Name: req.GetUsername(),
		Role: 2,
	},
		refreshTokenSecretKey,
		refreshTokenExpiration,
	)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &auth.LoginResponse{RefreshToken: refreshToken}, nil
}

func (s *serverAuth) GetRefreshToken(ctx context.Context, req *auth.GetRefreshTokenRequest) (*auth.GetRefreshTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetOldRefreshToken(), refreshTokenSecretKey)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "incalid refresh token")
	}

	refreshToken, err := utils.GenerateToken(model.User{
		Name: claims.Name,
		//в реальности лезем и берем данные из кеша или бд
		Role: 2,
	},
		refreshTokenSecretKey,
		refreshTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	return &auth.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}

func (s *serverAuth) GetAccessToken(ctx context.Context, req *auth.GetAccessTokenRequest) (*auth.GetAccessTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), refreshTokenSecretKey)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "incalid refresh token")
	}

	accessToken, err := utils.GenerateToken(model.User{
		Name: claims.Name,
		//в реальности лезем и берем данные из кеша или бд
		Role: 2,
	},
		accessTokenSecretKey,
		accessTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	return &auth.GetAccessTokenResponse{AccessToken: accessToken}, nil
}

func (s *serverAccess) Check(ctx context.Context, req *access.CheckRequest) (*emptypb.Empty, error) {
	//из ctx вынимаем метаданные
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	//получаем заголовок authorization
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	//извлекаем сам токен
	accessToken := strings.TrimSpace(strings.TrimPrefix(authHeader[0], authPrefix))

	//проверка токена
	claims, err := utils.VerifyToken(accessToken, accessTokenSecretKey)
	if err != nil {
		log.Printf("Token verification failed: %s", err.Error()) // Логируем ошибку проверки токена
		return nil, errors.New("access token is invalid")
	}

	//получение ролей доступных из эндпоинтов
	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		return nil, errors.New("failed to get accessible roles")
	}

	//проверка доступа
	role, ok := accessibleMap[req.GetEndpointAddress()]
	if !ok {
		return &emptypb.Empty{}, nil
	}

	//сравнение роли из токена с требуемой ролью
	if role == claims.Role {
		return &emptypb.Empty{}, nil
	}

	//отказ в доступе
	return nil, errors.New("access denied")
}

func (s *serverAccess) accessibleRoles(ctx context.Context) (map[string]model.UserRole, error) {
	if accessRoles == nil {
		accessRoles = make(map[string]model.UserRole)

		//лезем в бд за данными о доступных ролях для каждого эндпоинта
		//можно кэшировать данные, чтобы не лезть в бд каждый раз

		//например для эндпоинта /user/get доступна только роль admin
		accessRoles[model.ExamplePath] = model.UserRoleAdmin
	}

	return accessRoles, nil
}

var accessRoles map[string]model.UserRole

var err error

func runPrometheus() error {
	log.Println("Starting Prometheus metrics server...")

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:    "localhost:2112",
		Handler: mux,
	}

	log.Printf("Prometheus server is running on %s", prometheusServer.Addr)

	err := prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

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

	userRepo := userRepoPkg.NewRepository(db)
	userServ := userServPkg.NewService(userRepo, cache)

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
		grpc.UnaryInterceptor(
			interceptor.MetricsInterceptor,
		),
	) /*grpc.Creds(creds) - тут tls соединение*/
	pb.RegisterUserApiServer(grpcServer, userAPIPkg.NewImplementation(userServ))

	lis, err = net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052", err)
	}

	go func() {
		log.Printf("gRPC server is running on port 50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	go func() {
		err = runPrometheus()
		if err != nil {
			log.Fatal(err)
		}
	}()

	//тестим auth
	grpcServerAuth := grpc.NewServer()
	auth.RegisterAuthApiServer(grpcServerAuth, &serverAuth{})

	lisAuth, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen on port 50053", err)
	}

	go func() {
		log.Printf("gRPC Auth server is running on port 50053")
		if err := grpcServerAuth.Serve(lisAuth); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	//тестим access
	grpcServerAccess := grpc.NewServer()
	access.RegisterAccessServer(grpcServerAccess, &serverAccess{})

	lisAccess, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Failed to listen on port 50054: %v", err)
	}

	go func() {
		log.Printf("gRPC Access server is running on port 50054")
		if err := grpcServerAccess.Serve(lisAccess); err != nil {
			log.Fatalf("Failed to serve Access gRPC: %v", err)
		}
	}()

	select {}

}
