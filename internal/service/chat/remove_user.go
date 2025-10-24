package chat

import (
	"context"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) RemoveUser(ctx context.Context, senderID, chatID, userID int64) error {
	targetRole, err := s.checkUserRole(ctx, chatID, userID, model.ROLE_USER)
	if err != nil {
		return err
	}

	senderRole, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_ADMIN)
	if err != nil {
		return err
	}

	if senderRole <= targetRole {
		return domainerrors.ErrCannotRemoveHigherUser
	}

	if err := s.chatParticipantRepo.RemoveUser(ctx, chatID, userID); err != nil {
		return err
	}

	s.streams.RemoveMsgStream(chatID, userID)

	return nil
}
