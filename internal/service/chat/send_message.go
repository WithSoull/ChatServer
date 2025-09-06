package chat

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) SendMessage(ctx context.Context, senderID int64, chatID int64, text string) error {
	if text == "" {
		return status.Errorf(codes.InvalidArgument, "message text cannot be empty")
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
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("failed to create message in chat %d from sender %d: %v", chatID, senderID, err)
		}
		return grpcErr
	}

	return nil
}
