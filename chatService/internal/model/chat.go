package model

import (
	"time"
)

// Чат
type Chat struct {
	ID        int64
	CreatedAt time.Time
	Usernames []string // Список пользователей в чате
}

// Сообщение
type Message struct {
	ID        int64
	ChatID    int64
	FromUser  string
	Text      string
	Timestamp time.Time
}

type ChatUser struct {
	ChatID   int64
	Username string
}
