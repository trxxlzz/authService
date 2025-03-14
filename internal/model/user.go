package model

import "time"

type UserRole int

const (
	UserRoleUnspecified UserRole = 0
	UserRoleUser        UserRole = 1
	UserRoleAdmin       UserRole = 2
)

type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      UserRole  `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
