package chat

import (
	"context"
	"errors"
)

func (s *Service) EditMessage(
	ctx context.Context,
	senderID, messageID int64,
	newText *string,
	newIsPinned *bool,
) error {
	msg, err := s.msgRepo.Get(ctx, messageID)
	if err != nil {
		return err
	}

	ok, senderRole := s.chatRepo.GetUserRole(ctx, msg.ChatID, senderID)
	if !ok {
		return errors.New("you are not a member of the chat")
	}

	if newText == nil && newIsPinned == nil {
		return errors.New("no changes provided")
	}

	if newText != nil {
		if senderID != msg.SenderID {
			return errors.New("you can edit only your own messages")
		}
		msg.Text = *newText
	}

	if newIsPinned != nil {
		if senderRole < 1 {
			return errors.New("only admins or owner can pin/unpin messages")
		}
		msg.IsPinned = *newIsPinned
	}

	return s.msgRepo.Update(ctx, msg)
}
