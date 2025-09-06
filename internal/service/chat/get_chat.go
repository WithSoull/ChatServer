package chat

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) GetChat(ctx context.Context, senderID, chatID int64) (model.Chat, []model.Message, error) {
	if _, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_USER); err != nil {
		return model.Chat{}, nil, err
	}

	chat, err := s.chatRepo.Get(ctx, chatID)
	if err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to get chat: %v", err)
		}
		return model.Chat{}, nil, grpcErr
	}

	chatParticipants, err := s.chatParticipantRepo.GetUsers(ctx, chatID)
	if err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to get users of the chat: %v", err)
		}
		return model.Chat{}, nil, grpcErr
	}
	chat.ChatInfo.UserIDs = chatParticipants

	msgs, err := s.msgRepo.GetByChat(ctx, chatID)
	if err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to get messages of the chat: %v", err)
		}
		return model.Chat{}, nil, grpcErr
	}

	return chat, msgs, nil
}
