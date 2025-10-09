package chat

import (
	"context"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) EditMessage(ctx context.Context, senderID, messageID int64, newText string) error {
	// Validation
	if newText == "" {
		return domainerrors.ErrNoChangesProvided
	}

	msg, err := s.msgRepo.Get(ctx, messageID)
	if err != nil {
		return err
	}

	if senderID != msg.SenderID {
		return domainerrors.ErrCannotEditOthersMessages
	}

	if _, err := s.checkUserRole(ctx, msg.ChatID, senderID, model.ROLE_USER); err != nil {
		return err
	}

	msg.Text = newText

	if err := s.msgRepo.Update(ctx, msg); err != nil {
		return err
	}

	return nil
}
