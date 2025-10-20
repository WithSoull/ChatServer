package app

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/config"
	chatHandler "github.com/WithSoull/ChatServer/internal/handler/chat"
	"github.com/WithSoull/ChatServer/internal/repository"
	chatRepository "github.com/WithSoull/ChatServer/internal/repository/chat"
	chatParticipantRepository "github.com/WithSoull/ChatServer/internal/repository/chat_participant"
	msgRepository "github.com/WithSoull/ChatServer/internal/repository/message"
	"github.com/WithSoull/ChatServer/internal/service"
	chatService "github.com/WithSoull/ChatServer/internal/service/chat"
	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"github.com/WithSoull/platform_common/pkg/client/db"
	"github.com/WithSoull/platform_common/pkg/client/db/pg"
	"github.com/WithSoull/platform_common/pkg/client/db/transaction"
	"github.com/WithSoull/platform_common/pkg/closer"
	"github.com/WithSoull/platform_common/pkg/logger"
)

type serviceProvider struct {
	pgClient  db.Client
	txManager db.TxManager

	chatRepo            repository.ChatRepo
	chatParticipantRepo repository.ChatParticipantRepo
	msgRepo             repository.MessageRepo

	chatService service.ChatService
	chatHandler desc.ChatV1Server
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGClient(ctx context.Context) db.Client {
	if s.pgClient == nil {
		client, err := pg.NewPGClient(ctx, logger.Logger(), config.AppConfig().PG)
		if err != nil {
			panic(err)
		}

		if err := client.DB().Ping(ctx); err != nil {
			panic(err)
		}

		closer.AddNamed("PGClient", func(ctx context.Context) error {
			return client.Close()
		})

		s.pgClient = client
	}

	return s.pgClient
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepo {
	if s.chatRepo == nil {
		s.chatRepo = chatRepository.NewRepository(s.PGClient(ctx))
	}

	return s.chatRepo
}

func (s *serviceProvider) MessageRepository(ctx context.Context) repository.MessageRepo {
	if s.msgRepo == nil {
		s.msgRepo = msgRepository.NewRepository(s.PGClient(ctx))
	}

	return s.msgRepo
}

func (s *serviceProvider) ChatParticipantRepository(ctx context.Context) repository.ChatParticipantRepo {
	if s.chatParticipantRepo == nil {
		s.chatParticipantRepo = chatParticipantRepository.NewRepository(s.PGClient(ctx))
	}

	return s.chatParticipantRepo
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.PGClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(s.ChatRepository(ctx), s.MessageRepository(ctx), s.ChatParticipantRepository(ctx), s.TxManager(ctx))
	}

	return s.chatService
}

func (s *serviceProvider) ChatHandler(ctx context.Context) desc.ChatV1Server {
	if s.chatHandler == nil {
		s.chatHandler = chatHandler.NewHandler(s.ChatService(ctx))
	}

	return s.chatHandler
}
