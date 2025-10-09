package chat

import (
	"context"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) CreateChat(ctx context.Context, senderID int64, chat_info model.ChatInfo) (int64, error) {
	var chatID int64

	// Validation
	userSet := make(map[int64]struct{})
	for _, id := range chat_info.UserIDs {
		if _, exists := userSet[id]; exists {
			return 0, domainerrors.ErrDuplicateParticipant
		}
		userSet[id] = struct{}{}
	}

	txErr := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		chat := model.Chat{
			OwnerID:  senderID,
			ChatInfo: chat_info,
		}
		var err error
		chatID, err = s.chatRepo.Create(ctx, chat)
		if err != nil {
			return err
		}

		isOwnerAdded := false
		for userID, _ := range userSet {
			var userRole model.Role
			if senderID == userID {
				isOwnerAdded = true
				userRole = model.ROLE_OWNER
			} else {
				userRole = model.ROLE_USER
			}

			if err := s.chatParticipantRepo.AddUser(ctx, chatID, userID, userRole); err != nil {
				return err
			}
		}

		if !isOwnerAdded {
			err := s.chatParticipantRepo.AddUser(ctx, chatID, senderID, model.ROLE_OWNER)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if txErr != nil {
		return 0, txErr
	}
	return chatID, nil
}
