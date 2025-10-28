package user

import (
	"context"

	"github.com/WithSoull/platform_common/pkg/kafka"
	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
)

func (s *userConcumerService) UserCreatedHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.userCreatedDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode UserCreated", zap.Error(err))
		return err
	}

	s.cache.Add(ctx, event.UserID)

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.Int64("userID", event.UserID),
		zap.Timep("createdAt", event.CreatedAt),
	)

	return nil
}

func (s *userConcumerService) UserDeletedHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.userDeletedDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode UserDeleted", zap.Error(err))
		return err
	}

	s.cache.Remove(ctx, event.UserID)
	err = s.chatService.RemoveUserFromAllChats(ctx, event.UserID)
	if err != nil {
		logger.Error(ctx, "failed to remove deleted user from all chats", zap.Error(err), zap.Int64("userID", event.UserID))
		return err
	}

	logger.Info(
		ctx,
		"Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.Int64("userID", event.UserID),
		zap.Timep("deletedAt", event.DeletedAt),
	)

	return nil
}
