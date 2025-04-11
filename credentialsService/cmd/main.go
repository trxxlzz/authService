package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"microservices/credentialsService/internal/kafka"
	"microservices/credentialsService/internal/service/role"
	"microservices/credentialsService/internal/transport/rest"
	"net/http"
)

func main() {
	// Инициализация Kafka writer
	writer := kafka.NewKafkaWriter("localhost:9092", "add-role")
	defer writer.Close()

	// Инициализация бизнес-слоя
	kafkaService := role.NewKafkaService(writer)

	// Создание роутера
	r := chi.NewRouter()
	r.Post("/add-role", rest.RoleHandler(kafkaService))

	// Запуск HTTP сервера на порту 8081
	log.Println("HTTP server listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
