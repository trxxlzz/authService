package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// NewMongoConnection подключается к MongoDB
func NewMongoConnection(ctx context.Context, uri string) (*mongo.Client, error) {
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to MongoDB: %v", err)
	}

	// Проверяем подключение
	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctxPing, nil); err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %v", err)
	}

	return client, nil
}
