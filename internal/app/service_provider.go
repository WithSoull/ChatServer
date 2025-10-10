package app

import (
	"context"
	"log"

	"github.com/WithSoull/ChatServer/internal/config"
	"github.com/WithSoull/ChatServer/internal/config/env"
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
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig
	httpConfig config.HTTPConfig

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

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) PGClient(ctx context.Context) db.Client {
	if s.pgClient == nil {
		client, err := pg.NewPGClient(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create connection pool: %s", err.Error())
		}

		if err := client.DB().Ping(ctx); err != nil {
			log.Fatalf("failed to connect to database: %v", err.Error())
		}

		closer.Add(func() error {
			client.Close()
			return nil
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
