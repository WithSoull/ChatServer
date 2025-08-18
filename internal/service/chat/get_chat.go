package chat

import (
	"context"
	"errors"

	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) GetChat(ctx context.Context, senderID, chatID int64) (model.Chat, []model.Message, error) {
	if ok, _ := s.chatRepo.GetUserRole(ctx, chatID, senderID); ok {
		chat, err := s.chatRepo.Get(ctx, chatID)
		if err != nil {
			return model.Chat{}, nil, err
		}

		msgs, err := s.msgRepo.GetByChat(ctx, chatID)
		if err != nil {
			return model.Chat{}, nil, err
		}

		return chat, msgs, nil
	}

	return model.Chat{}, nil, errors.New("You dont exist in this chat")
}
