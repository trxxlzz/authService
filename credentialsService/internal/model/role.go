package model

type RoleMessage struct {
	UserID string `json:"user_id"`
	Role   int    `json:"role"`
}
