package chat

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) DeleteChat(ctx context.Context, senderID, chatID int64) error {
	if _, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_OWNER); err != nil {
		return err
	}

	if err := s.chatRepo.Delete(ctx, chatID); err != nil {
		return err
	}

	s.streams.RemoveChatStream(chatID)
	return nil
}
