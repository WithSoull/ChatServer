package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/WithSoull/ChatServer/internal/client/cache"
	"github.com/WithSoull/ChatServer/internal/client/cache/redis"
	"github.com/WithSoull/ChatServer/internal/config"
	converterKafka "github.com/WithSoull/ChatServer/internal/converter/kafka"
	decoderKafka "github.com/WithSoull/ChatServer/internal/converter/kafka/decoder"
	chatHandler "github.com/WithSoull/ChatServer/internal/handler/chat"
	"github.com/WithSoull/ChatServer/internal/repository"
	chatRepository "github.com/WithSoull/ChatServer/internal/repository/chat"
	chatParticipantRepository "github.com/WithSoull/ChatServer/internal/repository/chat_participant"
	msgRepository "github.com/WithSoull/ChatServer/internal/repository/message"
	"github.com/WithSoull/ChatServer/internal/service"
	chatService "github.com/WithSoull/ChatServer/internal/service/chat"
	userConsumerService "github.com/WithSoull/ChatServer/internal/service/consumer/user"
	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"github.com/WithSoull/platform_common/pkg/client/db"
	"github.com/WithSoull/platform_common/pkg/client/db/pg"
	"github.com/WithSoull/platform_common/pkg/client/db/transaction"
	"github.com/WithSoull/platform_common/pkg/closer"
	platformKafka "github.com/WithSoull/platform_common/pkg/kafka"
	platformKafkaConsumer "github.com/WithSoull/platform_common/pkg/kafka/consumer"
	"github.com/WithSoull/platform_common/pkg/logger"
	"github.com/WithSoull/platform_common/pkg/tokens"
	"github.com/WithSoull/platform_common/pkg/tokens/jwt"
	redigo "github.com/gomodule/redigo/redis"
)

type serviceProvider struct {
	pgClient  db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	cacheClient cache.UsersIDsCacheClient

	chatRepo            repository.ChatRepo
	chatParticipantRepo repository.ChatParticipantRepo
	msgRepo             repository.MessageRepo

	chatService   service.ChatService
	tokenVerifier tokens.TokenVerifier
	chatHandler   desc.ChatV1Server

	userCreatedConsumerGroup sarama.ConsumerGroup
	userCreatedConsumer      platformKafka.Consumer
	userCreatedDecoder       converterKafka.UserCreatedDecoder

	userDeletedConsumerGroup sarama.ConsumerGroup
	userDeletedConsumer      platformKafka.Consumer
	userDeletedDecoder       converterKafka.UserDeletedDecoder

	userConsumerService service.UserConsumerService
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

func (s *serviceProvider) CacheClient(ctx context.Context) cache.UsersIDsCacheClient {
	if s.cacheClient == nil {
		client := redis.NewClient(s.RedisPool())
		s.cacheClient = client
	}

	return s.cacheClient
}

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		redisPool := &redigo.Pool{
			MaxIdle:     int(config.AppConfig().Redis.MaxIdle()),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}

		closer.AddNamed("RedisPool", func(ctx context.Context) error {
			return redisPool.Close()
		})

		s.redisPool = redisPool
	}

	return s.redisPool
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
		s.chatService = chatService.NewService(
			s.ChatRepository(ctx),
			s.MessageRepository(ctx),
			s.ChatParticipantRepository(ctx),
			s.TxManager(ctx),
			s.CacheClient(ctx),
		)
	}

	return s.chatService
}

func (s *serviceProvider) TokenVerifier(ctx context.Context) tokens.TokenVerifier {
	if s.tokenVerifier == nil {
		s.tokenVerifier = jwt.NewJWTVerifier(config.AppConfig().JWT)
	}
	return s.tokenVerifier
}

func (s *serviceProvider) ChatHandler(ctx context.Context) desc.ChatV1Server {
	if s.chatHandler == nil {
		s.chatHandler = chatHandler.NewHandler(s.ChatService(ctx), s.TokenVerifier(ctx))
	}

	return s.chatHandler
}

func (s *serviceProvider) UserCreatedConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	if s.userCreatedConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().Sarama.GroupID(),
			config.AppConfig().Sarama.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %v\n", err))
		}

		closer.AddNamed("Kafka user.created consumer group", func(ctx context.Context) error {
			return s.userCreatedConsumerGroup.Close()
		})

		s.userCreatedConsumerGroup = consumerGroup
	}

	return s.userCreatedConsumerGroup
}

func (s *serviceProvider) UserDeletedConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	if s.userDeletedConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().Sarama.GroupID(),
			config.AppConfig().Sarama.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create user.deleted consumer group: %v\n", err))
		}

		closer.AddNamed("Kafka user.deleted consumer group", func(ctx context.Context) error {
			return s.userDeletedConsumerGroup.Close()
		})

		s.userDeletedConsumerGroup = consumerGroup
	}

	return s.userDeletedConsumerGroup
}

func (s *serviceProvider) UserCreatedConsumer(ctx context.Context) platformKafka.Consumer {
	if s.userCreatedConsumer == nil {
		s.userCreatedConsumer = platformKafkaConsumer.NewConsumer(
			s.UserCreatedConsumerGroup(ctx),
			[]string{
				config.AppConfig().UserCreatedConsumer.Topic(),
			},
			logger.Logger(),
		)
	}

	return s.userCreatedConsumer
}

func (s *serviceProvider) UserCreatedDecoder(ctx context.Context) converterKafka.UserCreatedDecoder {
	if s.userCreatedDecoder == nil {
		s.userCreatedDecoder = decoderKafka.NewUserCreatedDecoder()
	}
	return s.userCreatedDecoder
}

func (s *serviceProvider) UserDeletedConsumer(ctx context.Context) platformKafka.Consumer {
	if s.userDeletedConsumer == nil {
		s.userDeletedConsumer = platformKafkaConsumer.NewConsumer(
			s.UserDeletedConsumerGroup(ctx),
			[]string{
				config.AppConfig().UserDeletedConsumer.Topic(),
			},
			logger.Logger(),
		)
	}

	return s.userDeletedConsumer
}

func (s *serviceProvider) UserDeletedDecoder(ctx context.Context) converterKafka.UserDeletedDecoder {
	if s.userDeletedDecoder == nil {
		s.userDeletedDecoder = decoderKafka.NewUserDeletedDecoder()
	}
	return s.userDeletedDecoder
}

func (s *serviceProvider) UserConsumerService(ctx context.Context) service.UserConsumerService {
	if s.userConsumerService == nil {
		s.userConsumerService = userConsumerService.NewService(
			s.UserCreatedConsumer(ctx),
			s.UserCreatedDecoder(ctx),
			s.UserDeletedConsumer(ctx),
			s.UserDeletedDecoder(ctx),
			s.CacheClient(ctx),
			s.ChatService(ctx),
		)
	}
	return s.userConsumerService
}
