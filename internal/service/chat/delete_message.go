package chat

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) DeleteMessage(ctx context.Context, senderID, messageID int64) error {
	msg, err := s.msgRepo.Get(ctx, messageID)
	if err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("failed to get sender for message %d: %v", messageID, err)
		}
		return grpcErr
	}

	if msg.SenderID != senderID {
		messageSenderRole, err := s.checkUserRole(ctx, msg.ChatID, msg.SenderID, model.ROLE_USER)
		if err != nil {
			return err
		}

		// We need to check that the role of request's sender is higher than
		// the role of message's sender, thats why we increment neededRole while checking
		if _, err := s.checkUserRole(ctx, msg.ChatID, senderID, messageSenderRole+1); err != nil {
			return err
		}
	}

	if err := s.msgRepo.Delete(ctx, messageID); err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("failed to delete message %d: %v", messageID, err)
		}
		return grpcErr
	}

	return nil
}
