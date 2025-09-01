package rmodel

import "time"

type Role int64

const (
	ROLE_USER Role = iota
	ROLE_ADMIN
	ROLE_OWNER
)

// Table chats
type Chat struct {
	ID          int64     `db:"id"`
	OwnerID     int64     `db:"owner_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Table chat_participants
type ChatParticipant struct {
	ChatID    int64     `db:"chat_id"`
	UserID    int64     `db:"user_id"`
	Role      Role      `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Table messages
type Message struct {
	ID        int64     `db:"id"`
	ChatID    int64     `db:"chat_id"`
	SenderID  int64     `db:"sender_id"`
	Text      string    `db:"text"`
	IsPinned  bool      `db:"is_pinned"`
	SendAt    time.Time `db:"send_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
