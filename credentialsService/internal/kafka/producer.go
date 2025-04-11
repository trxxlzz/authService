package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"microservices/credentialsService/internal/model"
	"time"
)

type Writer interface {
	WriteMessage(msg model.RoleMessage) error
	Close() error
}

type KafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaWriter(broker, topic string) *KafkaWriter {
	w := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaWriter{writer: w}
}

func (k *KafkaWriter) WriteMessage(msg model.RoleMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return k.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(msg.UserID),
		Value: data,
		Time:  time.Now(),
	})
}

func (k *KafkaWriter) Close() error {
	return k.writer.Close()
}
