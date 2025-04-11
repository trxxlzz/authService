package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"microservices/authService/internal/mapper"
	"microservices/authService/internal/service"
	"time"
)

type Consumer struct {
	reader      *kafka.Reader
	userService service.UserService
}

func NewConsumer(brokers []string, topic string, userService service.UserService) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
			MaxWait:  time.Second,
		}),
		userService: userService,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	log.Println("Kafka consumer started")
	defer log.Println("Kafka consumer stopped")

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer stopped")
			return nil
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Kafka read error: %v", err)
				continue
			}

			// Используем маппер
			id, role, err := mapper.ToUserFromKafka(&msg)
			if err != nil {
				log.Printf("Failed to map Kafka message: %v", err)
				continue
			}

			if err := c.userService.UpdateUserRole(ctx, id, role); err != nil {
				log.Printf("Failed to update user role: %v (UserID: %s, Role: %d)", err, id, role)
			} else {
				log.Printf("Updated role for UserID %s to %d", id, role)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
