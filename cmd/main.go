package main

import (
	"authService/internal/config"
	"authService/internal/infra"
	"context"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := infra.NewDBConnection(cfg.DSN())
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close(context.Background())

	log.Println("Successfully connected to database")
}
