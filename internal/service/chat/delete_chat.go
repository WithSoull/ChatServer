package chat

import (
	"context"
	"errors"

	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) DeleteChat(ctx context.Context, senderID, chatID int64) error {
	if ok, role := s.chatRepo.GetUserRole(ctx, chatID, senderID); ok && (role == model.ROLE_OWNER) {
		return s.chatRepo.Delete(ctx, chatID)
	}
	return errors.New("Only owner of the chat can delete it")
}
