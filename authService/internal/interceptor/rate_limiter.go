package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices/authService/internal/rate_limiter"
)

// RateLimiterInterceptor — обёртка для токен-бакета
type RateLimiterInterceptor struct {
	rateLimiter *rate_limiter.TokenBucketLimiter
}

// Конструктор: создаёт новый интерцептор с заданным rate limiter'ом
func NewRateLimiterInterceptor(rateLimiter *rate_limiter.TokenBucketLimiter) *RateLimiterInterceptor {
	return &RateLimiterInterceptor{
		rateLimiter: rateLimiter,
	}
}

// Unary — функция, вызываемая для каждого unary gRPC запроса
func (r *RateLimiterInterceptor) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Проверка: доступен ли токен?
	if !r.rateLimiter.Allow() {
		return nil, status.Error(codes.ResourceExhausted, "too many requests")
	}

	return handler(ctx, req)
}
