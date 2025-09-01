package chat_participant

import (
	"github.com/WithSoull/ChatServer/internal/client/db"
	"github.com/WithSoull/ChatServer/internal/repository"
)

const (
	TableName       = "chat_participants"
	ChatIDColumn    = "chat_id"
	UserIDColumn    = "user_id"
	RoleColumn      = "role"
	CreatedAtColumn = "joined_at"
	UpdatedAtColumn = "updated_at"
)

type chatParticipantRepo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.ChatParticipantRepo {
	return &chatParticipantRepo{db: db}
}
