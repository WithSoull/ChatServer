package chat

import (
	"context"
	"errors"
	"strings"

	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) SendMessage(ctx context.Context, fromUserID int64, chatID int64, text string) errror {
	// Example of censor obscene language
	if strings.Contains(text, "duck") {
		return errors.New("We dont love ducks bro")
	}

	msg := model.Message{
		FromUserID: fromUserID,
		ChatID:     chatID,
		Text:       text,
	}

	return s.repo.SendMessage(ctx, fromUserID, chatID, text)

}

