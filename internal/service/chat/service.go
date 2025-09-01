package chat

import (
	"github.com/WithSoull/ChatServer/internal/client/db"
	"github.com/WithSoull/ChatServer/internal/repository"
)

type Service struct {
	chatRepo            repository.ChatRepo
	msgRepo             repository.MessageRepo
	chatParticipantRepo repository.ChatParticipantRepo

	txManager db.TxManager
}

func NewService(
	chatRepo repository.ChatRepo,
	msgRepo repository.MessageRepo,
	chatParticipantRepo repository.ChatParticipantRepo,
	txManager db.TxManager,
) *Service {
	return &Service{
		chatRepo:            chatRepo,
		msgRepo:             msgRepo,
		chatParticipantRepo: chatParticipantRepo,
		txManager:           txManager,
	}
}
