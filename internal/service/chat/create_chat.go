package chat

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) CreateChat(ctx context.Context, senderID int64, chat_info model.ChatInfo) (int64, error) {
	chat := model.Chat{
		OwnerID:  senderID,
		ChatInfo: chat_info,
	}
	chat_id, err := s.chatRepo.Create(ctx, chat)
	return chat_id, err
}
