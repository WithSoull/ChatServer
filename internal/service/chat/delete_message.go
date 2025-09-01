package chat

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) DeleteMessage(ctx context.Context, senderID, messageID int64) error {
	msgSenderID, err := s.msgRepo.GetMessageSenderID(ctx, messageID)
	if err != nil {
		log.Printf("failed to get sender for message %d: %v", messageID, err)
		return status.Errorf(codes.NotFound, "message not found")
	}

	if msgSenderID != senderID {
		return status.Errorf(codes.PermissionDenied, "only the sender of a message can delete it")
	}

	if err := s.msgRepo.Delete(ctx, messageID); err != nil {
		log.Printf("failed to delete message %d: %v", messageID, err)
		return status.Errorf(codes.Internal, "failed to delete message")
	}

	return nil
}
