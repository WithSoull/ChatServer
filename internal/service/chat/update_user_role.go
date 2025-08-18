package chat

import (
	"context"
	"errors"

	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) UpdateUserRole(ctx context.Context, senderID, chatID, userID int64, newRole model.Role) error {
	ok, senderRole := s.chatRepo.GetUserRole(ctx, chatID, senderID)
	if !ok {
		return errors.New("you are not a member of the chat")
	}

	ok, targetRole := s.chatRepo.GetUserRole(ctx, chatID, userID)
	if !ok {
		return errors.New("target user is not a member of the chat")
	}

	if newRole > senderRole {
		return errors.New("cannot assign a role higher than your own")
	}

	if targetRole >= senderRole {
		return errors.New("cannot change role of a user with equal or higher role")
	}

	return s.chatRepo.UpdateUserRole(ctx, chatID, userID, newRole)
}
