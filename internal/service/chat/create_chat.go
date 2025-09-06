package chat

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) CreateChat(ctx context.Context, senderID int64, chat_info model.ChatInfo) (int64, error) {
	var chatID int64

	userSet := make(map[int64]struct{})
	for _, id := range chat_info.UserIDs {
		if _, exists := userSet[id]; exists {
			return 0, status.Error(codes.InvalidArgument, "Each participant can only be added once")
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
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(txErr)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to create chat: %v", txErr)
		}
		return 0, grpcErr
	}
	return chatID, nil
}
