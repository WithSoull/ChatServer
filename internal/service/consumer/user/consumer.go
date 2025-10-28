package user

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/client/cache"
	"github.com/WithSoull/ChatServer/internal/config"
	kafkaConverter "github.com/WithSoull/ChatServer/internal/converter/kafka"
	"github.com/WithSoull/ChatServer/internal/service"
	"github.com/WithSoull/platform_common/pkg/kafka"
	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type userConcumerService struct {
	userCreatedConsumer kafka.Consumer
	userCreatedDecoder  kafkaConverter.UserCreatedDecoder

	userDeletedConsumer kafka.Consumer
	userDeletedDecoder  kafkaConverter.UserDeletedDecoder

	cache       cache.UsersIDsCacheClient
	chatService service.ChatService
}

func NewService(
	userCreatedConsumer kafka.Consumer,
	userCreatedDecoder kafkaConverter.UserCreatedDecoder,

	userDeletedConsumer kafka.Consumer,
	userDeletedDecoder kafkaConverter.UserDeletedDecoder,

	cacheClient cache.UsersIDsCacheClient,
	chatService service.ChatService,
) *userConcumerService {
	return &userConcumerService{
		userCreatedConsumer: userCreatedConsumer,
		userCreatedDecoder:  userCreatedDecoder,

		userDeletedConsumer: userDeletedConsumer,
		userDeletedDecoder:  userDeletedDecoder,

		cache:       cacheClient,
		chatService: chatService,
	}
}

func (s *userConcumerService) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting UserConsumerService")

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		logger.Info(ctx, "Starting consume topic", zap.String("topic", config.AppConfig().UserCreatedConsumer.Topic()))
		return s.userCreatedConsumer.Consume(ctx, s.UserCreatedHandler)
	})

	g.Go(func() error {
		logger.Info(ctx, "Starting consume topic", zap.String("topic", config.AppConfig().UserDeletedConsumer.Topic()))
		return s.userDeletedConsumer.Consume(ctx, s.UserDeletedHandler)
	})

	return g.Wait()
}
