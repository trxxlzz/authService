package service

import (
	"microservices/credentialsService/internal/model"
)

// Интерфейс для записи сообщений в Kafka
type KafkaService interface {
	ProduceRoleAssignment(msg model.RoleMessage) error
}
