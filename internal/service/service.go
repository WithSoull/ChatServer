package service

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/model"
)

type ChatService interface {
	ConnectChat(ctx context.Context, chatID int64) (chan *model.Message, error)
	DisconnectChat(ctx context.Context, chatID int64) error
	CreateChat(ctx context.Context, chat_info model.ChatInfo) (int64, error)
	DeleteChat(ctx context.Context, chatID int64) error
	GetChat(ctx context.Context, chatID int64) (model.Chat, []model.Message, error)
	AddUser(ctx context.Context, chatID, userID int64, role model.Role) error
	RemoveUser(ctx context.Context, chatID, userID int64) error
	UpdateUserRole(ctx context.Context, chatID, userID int64, role model.Role) error
	SendMessage(ctx context.Context, chatID int64, text string) error
	DeleteMessage(ctx context.Context, messageID int64) error
	EditMessage(ctx context.Context, messageID int64, newText string) error
	PinMessage(ctx context.Context, messageID int64, newIsPinned bool) error

	// Stuff methods for concumers
	RemoveUserFromAllChats(ctx context.Context, userID int64) error
}

type UserConsumerService interface {
	RunConsumer(ctx context.Context) error
}
