package user

import (
	"authService/internal/model"
	"authService/internal/repository"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type mongoRepo struct {
	collection *mongo.Collection
}

func NewMongoRepository(client *mongo.Client, dbName, collectionName string) repository.UserRepository {
	return &mongoRepo{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

// CreateUser - сохранение пользователя в MongoDB
func (r *mongoRepo) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	// Генерируем ID вручную
	user.ID = time.Now().UnixNano() // Или использовать atomic счетчик

	_, err := r.collection.InsertOne(ctx, bson.M{
		"id":        user.ID,
		"name":      user.Name,
		"email":     user.Email,
		"password":  user.Password,
		"role":      user.Role,
		"createdAt": user.CreatedAt,
	})
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// GetUser - получение пользователя по ID
func (r *mongoRepo) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	var user User

	err := r.collection.FindOne(ctx, bson.M{"id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return ToUserFromRepo(&user), nil
}

// UpdateUser - обновление данных пользователя
func (r *mongoRepo) UpdateUser(ctx context.Context, id int64, name string, email string) error {
	// Обновляем пользователя по числовому ID
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"id": id}, // Ищем пользователя по int64 ID
		bson.M{"$set": bson.M{
			"name":      name,
			"email":     email,
			"updatedAt": time.Now(),
		}},
	)
	if err != nil {
		return err
	}

	// Проверяем, найден ли пользователь
	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// DeleteUser - удаление пользователя
func (r *mongoRepo) DeleteUser(ctx context.Context, userID int64) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"id": userID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
