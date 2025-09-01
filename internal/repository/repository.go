package repository

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/model"
)

type ChatRepo interface {
	Create(ctx context.Context, chat model.Chat) (int64, error)
	Delete(ctx context.Context, chatID int64) error
	Get(ctx context.Context, chatID int64) (model.Chat, error)
}

type ChatParticipantRepo interface {
	AddUser(ctx context.Context, chatID, userID int64, role model.Role) error
	RemoveUser(ctx context.Context, chatID, userID int64) error
	GetUsers(ctx context.Context, chatID int64) ([]int64, error)
	GetUserRole(ctx context.Context, chatID, userID int64) (bool, model.Role)
	UpdateUserRole(ctx context.Context, chatID, userID int64, role model.Role) error
}

type MessageRepo interface {
	Create(ctx context.Context, msg model.Message) (int64, error)
	Delete(ctx context.Context, messageID int64) error
	Update(ctx context.Context, msg model.Message) error
	Get(ctx context.Context, messageID int64) (model.Message, error)
	GetMessageSenderID(ctx context.Context, messageID int64) (int64, error)
	GetByChat(ctx context.Context, chatID int64) ([]model.Message, error)
}
