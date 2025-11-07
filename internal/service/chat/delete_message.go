package chat

import (
	"context"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
)

func (s *Service) DeleteMessage(ctx context.Context, messageID int64) error {
	senderID, ok := claimsctx.ExtractUserID(ctx)
	if !ok {
		return domainerrors.ErrFailedToVerify
	}

	msg, err := s.msgRepo.Get(ctx, messageID)
	if err != nil {
		return err
	}

	if msg.SenderID != senderID {
		messageSenderRole, err := s.checkUserRole(ctx, msg.ChatID, msg.SenderID, model.ROLE_USER)
		if err != nil {
			return err
		}

		// We need to check that the role of request's sender is higher than
		// the role of message's sender, thats why we increment neededRole while checking
		if _, err := s.checkUserRole(ctx, msg.ChatID, senderID, messageSenderRole+1); err != nil {
			return err
		}
	}

	if err := s.msgRepo.Delete(ctx, messageID); err != nil {
		return err
	}

	return nil
}
