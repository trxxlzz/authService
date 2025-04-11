package rest

import (
	"encoding/json"
	"log"
	"microservices/credentialsService/internal/model"
	"microservices/credentialsService/internal/service"
	"net/http"
)

// Хэндлер для обработки запросов на создание роли
func RoleHandler(service service.KafkaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Чтение данных из тела запроса
		var msg model.RoleMessage

		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		// Запись сообщения в Kafka через интерфейс
		err := service.ProduceRoleAssignment(msg)
		if err != nil {
			log.Println("Kafka write error:", err)
			http.Error(w, "failed to write to Kafka", http.StatusInternalServerError)
			return
		}

		// Ответ клиенту
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("message sent to Kafka"))
	}
}
