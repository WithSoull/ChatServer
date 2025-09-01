package message

import (
	"github.com/WithSoull/ChatServer/internal/client/db"
	"github.com/WithSoull/ChatServer/internal/repository"
)

const (
	TableName       = "messages"
	IDColumn        = "id"
	ChatIDColumn    = "chat_id"
	SenderIDColumn  = "sender_id"
	TextColumn      = "text"
	IsPinnedColumn  = "is_pinned"
	SendAtColumn    = "send_at"
	UpdatedAtColumn = "updated_at"
)

type messageRepo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.MessageRepo {
	return &messageRepo{db: db}
}
