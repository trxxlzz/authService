package main

import (
	"authService/internal/model"
	"authService/pkg/jwt/gen/go/access"
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	flag.Parse()

	ctx := context.Background()

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDM1MzgwODIsIm5hbWUiOiJVdCBkZXNlcnVudCIsInJvbGUiOjJ9.AAFZTjKrK_FBcWyU6YcPIwAWa6nMMSiwIhmN3i4i2kU"
	log.Printf("Sending Access Token: %s", accessToken)

	md := metadata.New(map[string]string{"Authorization": "Bearer " + accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:50054"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to dial GRPC client: %v", err)
	}
	defer conn.Close()

	client := access.NewAccessClient(conn)

	_, err = client.Check(ctx, &access.CheckRequest{
		EndpointAddress: model.ExamplePath, // Замените на реальное значение
	})
	if err != nil {
		log.Fatalf("Check failed: %v", err)
	}

	fmt.Println("Access granted")
}
