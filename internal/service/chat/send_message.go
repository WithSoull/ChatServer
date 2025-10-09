package chat

import (
	"context"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) SendMessage(ctx context.Context, senderID int64, chatID int64, text string) error {
	// Validation
	if text == "" {
		return domainerrors.ErrEmptyMessageText
	}

	_, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_USER)
	if err != nil {
		return err
	}

	msg := model.Message{
		SenderID: senderID,
		ChatID:   chatID,
		Text:     text,
	}

	_, err = s.msgRepo.Create(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}
