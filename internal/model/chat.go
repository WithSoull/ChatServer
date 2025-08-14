package model

import "time"

type Role int

const (
	ROLE_USER Role = iota
	ROLE_ADMIN
	ROLE_OWNER
)

type ChatInfo struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	UserIDs     []int64 `json:"user_ids"`
}

type Chat struct {
	ChatID   int64    `json:"chat_id"`
	OwnerID  int64    `json:"owner_id"`
	ChatInfo ChatInfo `json:"chat_info"`
}

type Message struct {
	MessageID int64     `json:"message_id"`
	SenderID  int64     `json:"sender_id"`
	ChatID    int64     `json:"chat_id"`
	Text      string    `json:"text"`
	IsPinned  bool      `json:"is_pinned"`
	SendAt    time.Time `json:"send_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
