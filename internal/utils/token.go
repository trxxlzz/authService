package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"authService/internal/model"
	"github.com/dgrijalva/jwt-go"
)

// Декодирование ключа из base64
func decodeKey(encodedKey string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedKey)
}

func GenerateToken(user model.User, encodedSecretKey string, duration time.Duration) (string, error) {
	secretKey, err := decodeKey(encodedSecretKey)
	if err != nil {
		log.Printf("Failed to decode secret key: %v", err)
		return "", err
	}

	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Name: user.Name,
		Role: user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func VerifyToken(tokenStr string, encodedSecretKey string) (*model.UserClaims, error) {
	secretKey, err := decodeKey(encodedSecretKey)
	if err != nil {
		log.Printf("Failed to decode secret key: %v", err)
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Проверяем, что алгоритм HMAC-SHA256
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		},
	)
	if err != nil {
		log.Printf("Token verification failed: %s", err.Error())
		return nil, fmt.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
