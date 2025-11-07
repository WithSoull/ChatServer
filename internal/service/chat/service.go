package chat

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/client/cache"
	"github.com/WithSoull/ChatServer/internal/config"
	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/repository"
	"github.com/WithSoull/ChatServer/internal/service"
	"github.com/WithSoull/ChatServer/internal/service/chat/stream"
	"github.com/WithSoull/platform_common/pkg/client/db"
)

type Service struct {
	chatRepo            repository.ChatRepo
	msgRepo             repository.MessageRepo
	chatParticipantRepo repository.ChatParticipantRepo

	cache cache.UsersIDsCacheClient

	streams stream.ChatStreams

	txManager db.TxManager
}

func NewService(
	chatRepo repository.ChatRepo,
	msgRepo repository.MessageRepo,
	chatParticipantRepo repository.ChatParticipantRepo,
	txManager db.TxManager,

	cache cache.UsersIDsCacheClient,
) service.ChatService {
	return &Service{
		chatRepo:            chatRepo,
		msgRepo:             msgRepo,
		chatParticipantRepo: chatParticipantRepo,
		txManager:           txManager,
		streams:             *stream.NewChatStreams(config.AppConfig().Streaming.BufferSize()),

		cache: cache,
	}
}

func (s *Service) userExist(ctx context.Context, userID int64) error {
	exist, err := s.cache.Exist(ctx, userID)
	if err != nil {
		return err
	}

	if !exist {
		return domainerrors.ErrUserNotFound(userID)
	}

	return nil
}
