package chat

import (
	"context"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
)

func (s *Service) ConnectChat(ctx context.Context, chatID int64) (chan *model.Message, error) {
	senderID, ok := claimsctx.ExtractUserID(ctx)
	if !ok {
		return nil, domainerrors.ErrFailedToVerify
	}
	logger.Debug(ctx, "successfully extarcted userID", zap.Int64("userID", senderID))

	if err := s.userExist(ctx, senderID); err != nil {
		return nil, err
	}

	_, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_USER)
	if err != nil {
		return nil, err
	}

	return s.streams.MsgStream(chatID, senderID), nil
}

// WARN: Dont check user permisions !!!
func (s *Service) DisconnectChat(ctx context.Context, chatID int64) error {
	senderID, ok := claimsctx.ExtractUserID(ctx)
	if !ok {
		return domainerrors.ErrFailedToVerify
	}
	s.streams.RemoveMsgStream(chatID, senderID)
	return nil
}
