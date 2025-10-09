package chat

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) GetChat(ctx context.Context, senderID, chatID int64) (model.Chat, []model.Message, error) {
	if _, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_USER); err != nil {
		return model.Chat{}, nil, err
	}

	chat, err := s.chatRepo.Get(ctx, chatID)
	if err != nil {
		return model.Chat{}, nil, err
	}

	chatParticipants, err := s.chatParticipantRepo.GetUsers(ctx, chatID)
	if err != nil {
		return model.Chat{}, nil, err
	}
	chat.ChatInfo.UserIDs = chatParticipants

	msgs, err := s.msgRepo.GetByChat(ctx, chatID)
	if err != nil {
		return model.Chat{}, nil, err
	}

	return chat, msgs, nil
}
