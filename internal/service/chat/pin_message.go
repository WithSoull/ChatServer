package chat

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) PinMessage(ctx context.Context, senderID, messageID int64, newIsPinned bool) error {
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
