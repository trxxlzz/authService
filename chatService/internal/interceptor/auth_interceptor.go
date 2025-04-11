package interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	access "microservices/authService/pkg/jwt/gen/go/access"
)

// Интерцептор для проверки токена и доступа пользователя
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("method:", info.FullMethod)
	// 1. Получаем метаданные (заголовки)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	// 2. Достаём токен из заголовка Authorization
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization header is not provided")
	}

	// 3. Проверяем токен через Auth API
	token := authHeader[0]
	isValid, err := checkAccessWithAuthService(token, info.FullMethod)
	if err != nil || !isValid {
		return nil, status.Errorf(codes.PermissionDenied, "access denied: %v", err)
	}

	log.Print("method:", info.FullMethod)

	// 4. Если всё ОК — передаём запрос дальше
	return handler(ctx, req)
}

// Функция для проверки токена и прав через Auth API
func checkAccessWithAuthService(token string, method string) (bool, error) {
	// Устанавливаем соединение с Auth API
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials())) // Поменяй на адрес auth-service
	if err != nil {
		return false, fmt.Errorf("failed to connect to auth service: %v", err)
	}
	defer conn.Close()

	// Создаём gRPC-клиента
	client := access.NewAccessApiClient(conn)

	// Устанавливаем таймаут на запрос
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 🔥 Добавляем токен в metadata для запроса в authService
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)

	// Делаем запрос в Auth API
	_, err = client.Check(ctx, &access.CheckRequest{
		EndpointAddress: method,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
