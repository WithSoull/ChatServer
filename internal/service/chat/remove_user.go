package chat

import (
	"context"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
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

	if err := s.chatParticipantRepo.RemoveUserFromChat(ctx, chatID, userID); err != nil {
		return err
	}

	s.streams.RemoveMsgStream(chatID, userID)

	return nil
}

func (s *Service) RemoveUserFromAllChats(ctx context.Context, userID int64) error {
	chatParticipants, err := s.chatParticipantRepo.RemoveUserFromAllChats(ctx, userID)
	if err != nil {
		return err
	}

	// case with len(chatIDs) == 0 is normal,
	// because user can not consist in any chat

	for _, chatParticipant := range chatParticipants {
		if chatParticipant.Role == model.ROLE_OWNER {
			logger.Info(ctx, "owner of the chat has been deleted; this chat will be deleted.", zap.Int64("owner", userID), zap.Int64("chatID", chatParticipant.ChatID))
			if err := s.chatRepo.Delete(ctx, chatParticipant.ChatID); err != nil {
				logger.Error(ctx, "failed to delete chat", zap.Int64("chatID", chatParticipant.ChatID))
			}
		}
	}

	// Closing streams
	for _, chatParticipant := range chatParticipants {
		s.streams.RemoveMsgStream(chatParticipant.ChatID, userID)
	}

	return nil
}
