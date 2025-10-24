package chat

import (
	"github.com/WithSoull/ChatServer/internal/config"
	"github.com/WithSoull/ChatServer/internal/repository"
	"github.com/WithSoull/ChatServer/internal/service/chat/stream"
	"github.com/WithSoull/platform_common/pkg/client/db"
)

type Service struct {
	chatRepo            repository.ChatRepo
	msgRepo             repository.MessageRepo
	chatParticipantRepo repository.ChatParticipantRepo

	streams stream.ChatStreams

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
		streams:             *stream.NewChatStreams(config.AppConfig().Streaming.BufferSize()),
	}
}
