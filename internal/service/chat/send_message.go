package chat

import (
	"context"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/ChatServer/internal/validator"
	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
)

func (s *Service) SendMessage(ctx context.Context, chatID int64, text string) error {
	senderID, ok := claimsctx.ExtractUserID(ctx)
	if !ok {
		return domainerrors.ErrFailedToVerify
	}

	// Validation
	if err := validate.Validate(
		ctx,
		validator.ValidateNotEmptyMessage(text),
	); err != nil {
		return err
	}

	_, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_USER)
	if err != nil {
		return err
	}

	msg := model.Message{
		SenderID: senderID,
		ChatID:   chatID,
		IsPinned: false,
		Text:     text,
	}

	err = s.msgRepo.Create(ctx, &msg)
	if err != nil {
		return err
	}

	s.streams.AddMsgToChatStream(chatID, &msg)

	return nil
}
