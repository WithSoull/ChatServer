package chat

import (
	"github.com/WithSoull/AuthService/internal/client/db"
	"github.com/WithSoull/ChatServer/internal/repository"
)

const (
	chatsTableName      = "chats"
	chatsChatIdColumn   = "chat_id"
	chatsUserIdColumn   = "user_id"
	chatsJoinedAtColumn = "joined_at"

	messagesTableName       = "messages"
	messagesIdColumn        = "id"
	messagesChatIdColumn    = "chat_id"
	messagesUserIdColumn    = "from_user_id"
	messagesTextColumn      = "text"
	messagesTimestampColumn = "timestamp"

	participantsTableName    = "participants"
	participantsIdColumn     = "id"
	participantsChatIdColumn = "chat_id"
	participantsUserIdColumn = "user_id"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}
