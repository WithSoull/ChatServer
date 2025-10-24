package chat

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
)

func (s *Service) ConnectChat(ctx context.Context, senderID, chatID int64) (chan *model.Message, error) {
	_, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_USER)
	if err != nil {
		return nil, err
	}

	return s.streams.MsgStream(chatID, senderID), nil
}

// WARN: Dont check user permisions !!!
func (s *Service) DisconnectChat(ctx context.Context, senderID, chatID int64) {
	logger.Debug(ctx, "handle disconnectChat", zap.Int64("chatID", chatID), zap.Int64("userID", senderID))
	s.streams.RemoveMsgStream(chatID, senderID)
}
