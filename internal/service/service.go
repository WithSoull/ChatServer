package service

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/model"
)

type ChatService interface {
	CreateChat(ctx context.Context, senderID int64, name, description string, user_ids []int64) (int64, error)
	DeleteChat(ctx context.Context, senderID, chatID int64) error
	GetChat(ctx context.Context, senderID, chatID int64) (model.Chat, []model.Message, error)
	AddUser(ctx context.Context, senderID, chatID, userID int64, role model.Role) error
	RemoveUser(ctx context.Context, senderID, chatID, userID int64) error
	UpdateUserRole(ctx context.Context, senderID, chatID, userID int64, role model.Role) error
	SendMessage(ctx context.Context, senderID, chatID int64, text string) error
	DeleteMessage(ctx context.Context, senderID, messageID int64) error
	EditMessage(ctx context.Context, senderID, messageID int64, newText *string, newIsPinned *bool) error
}
