package chat

import (
	"context"
	"errors"

	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) AddUser(ctx context.Context, senderID, chatID, userID int64, role model.Role) error {
	// role > 1 is equal admin or owner
	ok, senderRole := s.chatRepo.GetUserRole(ctx, chatID, senderID)
	if !ok {
		return errors.New("you are not a member of the chat")
	}

	if senderRole < 2 {
		return errors.New("Only adnmins and owner of the chat can add new users")
	}
	return s.chatRepo.AddUser(ctx, chatID, userID, role)
}
