package chat

import (
	"context"
	"log"

	"github.com/WithSoull/ChatServer/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) CreateChat(ctx context.Context, senderID int64, chat_info model.ChatInfo) (int64, error) {
	var chatID int64

	txErr := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		chat := model.Chat{
			OwnerID:  senderID,
			ChatInfo: chat_info,
		}
		var err error
		chatID, err = s.chatRepo.Create(ctx, chat)
		if err != nil {
			log.Printf("[Service Layer] failed to create chat: %v", err)
			return err
		}

		isOwnerAdded := false
		for _, userID := range chat_info.UserIDs {
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
		log.Printf("[Service Layer] failed to create chat: %v", txErr)
		return 0, status.Errorf(codes.Internal, "failed to create chat")
	}
	return chatID, nil
}
