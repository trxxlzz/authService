package model

import "time"

type UserRole int

const (
	UserRoleUnspecified UserRole = 0
	UserRoleUser        UserRole = 1
	UserRoleAdmin       UserRole = 2
)

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}
