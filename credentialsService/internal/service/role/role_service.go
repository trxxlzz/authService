package role

import (
	"log"
	"microservices/credentialsService/internal/kafka"
	"microservices/credentialsService/internal/model"
)

// Структура для работы с Kafka
type KafkaServiceImpl struct {
	writer kafka.Writer
}

// Реализация метода для записи в Kafka
func (s *KafkaServiceImpl) ProduceRoleAssignment(msg model.RoleMessage) error {
	err := s.writer.WriteMessage(msg)
	if err != nil {
		log.Println("Error writing to Kafka:", err)
		return err
	}
	return nil
}
