package chat

import (
	"github.com/WithSoull/ChatServer/internal/client/db"
	"github.com/WithSoull/ChatServer/internal/repository"
)

const (
	TableName         = "chats"
	IDColumn          = "id"
	OwnerIDColumn     = "owner_id"
	NameColumn        = "name"
	DescriptionColumn = "description"
	CreatedAtColumn   = "created_at"
	UpdatedAtColumn   = "updated_at"
)

type chatRepo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.ChatRepo {
	return &chatRepo{db: db}
}
