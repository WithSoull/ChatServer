package chat

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) DeleteChat(ctx context.Context, senderID, chatID int64) error {
	if _, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_OWNER); err != nil {
		return err
	}

	if err := s.chatRepo.Delete(ctx, chatID); err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to delete chat: %v", err)
		}
		return grpcErr
	}
	return nil
}
