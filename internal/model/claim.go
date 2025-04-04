package model

import "github.com/dgrijalva/jwt-go"

const (
	ExamplePath = "/user/get"
)

type UserClaims struct {
	jwt.StandardClaims
	Name string   `json:"name"`
	Role UserRole `json:"role"`
}
