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
	UpdateUserRole(ctx context.Context, chatID, userID int64, role model.Role) error

	// GetUserRole cant define who from user and chat does not exist
	// that's why you should use checkUserRole from service layer to
	// define this 2 corner cases
	// To define the difference you need to select in chats table,
	// so chat_participant table does not have access to another table
	GetUserRole(ctx context.Context, chatID, userID int64) (model.Role, error)
}

type MessageRepo interface {
	Create(ctx context.Context, msg *model.Message) error
	Delete(ctx context.Context, messageID int64) error
	Update(ctx context.Context, msg model.Message) error
	Get(ctx context.Context, messageID int64) (model.Message, error)
	GetByChat(ctx context.Context, chatID int64) ([]model.Message, error)
}
