package chat

import (
	"context"
	"errors"
)

func (s *Service) RemoveUser(ctx context.Context, senderID, chatID, userID int64) error {
	// role > 1 is equal admin or owner
	if ok, sender_role := s.chatRepo.GetUserRole(ctx, chatID, senderID); ok && (sender_role > 1) {
		return s.chatRepo.RemoveUser(ctx, chatID, userID)
	}

	return errors.New("Only admins and owner of the chat can remove users")
}
