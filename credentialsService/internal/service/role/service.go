package role

import "microservices/credentialsService/internal/kafka"

func NewKafkaService(writer kafka.Writer) *KafkaServiceImpl {
	return &KafkaServiceImpl{writer: writer}
}
