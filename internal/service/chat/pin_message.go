package chat

import (
	"context"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
)

func (s *Service) PinMessage(ctx context.Context, messageID int64, newIsPinned bool) error {
	senderID, ok := claimsctx.ExtractUserID(ctx)
	if !ok {
		return domainerrors.ErrFailedToVerify
	}

	msg, err := s.msgRepo.Get(ctx, messageID)
	if err != nil {
		return err
	}

	if _, err := s.checkUserRole(ctx, msg.ChatID, senderID, model.ROLE_ADMIN); err != nil {
		return err
	}

	msg.IsPinned = newIsPinned

	if err := s.msgRepo.Update(ctx, msg); err != nil {
		return err
	}

	return nil
}
