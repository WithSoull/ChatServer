package chat

import (
	"context"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
)

func (s *Service) DeleteChat(ctx context.Context, chatID int64) error {
	senderID, ok := claimsctx.ExtractUserID(ctx)
	if !ok {
		return domainerrors.ErrFailedToVerify
	}

	if _, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_OWNER); err != nil {
		return err
	}

	if err := s.chatRepo.Delete(ctx, chatID); err != nil {
		return err
	}

	s.streams.RemoveChatStream(chatID)
	return nil
}
