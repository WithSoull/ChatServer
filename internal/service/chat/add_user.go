package chat

import (
	"context"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) AddUser(ctx context.Context, senderID, chatID, userID int64, role model.Role) error {
	// role > 1 is equal admin or owner
	senderRole, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_ADMIN)
	if err != nil {
		return err
	}
	if role >= senderRole {
		return domainerrors.ErrCannotAssignHigherRole
	}

	err = s.chatParticipantRepo.AddUser(ctx, chatID, userID, role)
	if err != nil {
		return err
	}

	return nil
}
