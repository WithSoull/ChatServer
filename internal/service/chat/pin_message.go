package chat

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) PinMessage(ctx context.Context, senderID, messageID int64, newIsPinned bool) error {
	msg, err := s.msgRepo.Get(ctx, messageID)
	if err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("failed to get message %d: %v", messageID, err)
		}
		return grpcErr
	}

	if _, err := s.checkUserRole(ctx, msg.ChatID, senderID, model.ROLE_ADMIN); err != nil {
		return err
	}

	msg.IsPinned = newIsPinned

	if err := s.msgRepo.Update(ctx, msg); err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("failed to update message %d: %v", messageID, err)
		}
		return grpcErr
	}

	return nil
}
