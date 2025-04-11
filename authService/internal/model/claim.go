package model

import "github.com/dgrijalva/jwt-go"

const (
	CreateChat  = "/proto.ChatApi/CreateChat"
	DeleteChat  = "/proto.ChatApi/DeleteChat"
	SendMessage = "/proto.ChatApi/SendMessage"
)

type UserClaims struct {
	jwt.StandardClaims
	Name string   `json:"name"`
	Role UserRole `json:"role"`
}
