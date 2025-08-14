package repository

import "github.com/WithSoull/ChatServer/internal/model"

type ChatRepo interface {
	Create(chat model.Chat) (int64, error)
	Delete(chatID int64) error
	Get(chatID int64) (model.Chat, error)
	AddUser(chatID, userID int64, role model.Role) error
	RemoveUser(chatID, userID int64) error
	UpdateUserRole(chatID, userID int64, role model.Role) error
}

type MessageRepo interface {
	Create(msg model.Message) error
	Delete(messageID int64) error
	Update(messageID int64, text *string, isPinned *bool) error
	GetByChat(chatID int64) ([]model.Message, error)
	Get(messageID int64) (model.Message, error)
}
