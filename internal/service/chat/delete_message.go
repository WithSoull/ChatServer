package chat

import (
	"context"
	"errors"
)

func (s *Service) DeleteMessage(ctx context.Context, senderID, messageID int64) error {
	msgSenderID := s.msgRepo.GetMessageSenderID(ctx, messageID)
	if msgSenderID != senderID {
		return errors.New("Only sender of message can delete it")
	}

	return s.msgRepo.Delete(ctx, messageID)
}
