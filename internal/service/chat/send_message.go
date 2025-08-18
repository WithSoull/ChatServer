package chat

import (
	"context"
	"errors"

	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) SendMessage(ctx context.Context, senderID int64, chatID int64, text string) error {
	ok, _ := s.chatRepo.GetUserRole(ctx, chatID, senderID)
	if !ok {
		return errors.New("You must be a member of the chat to send message")
	}

	msg := model.Message{
		SenderID: senderID,
		ChatID:   chatID,
		Text:     text,
	}

	return s.msgRepo.Create(ctx, msg)
}
