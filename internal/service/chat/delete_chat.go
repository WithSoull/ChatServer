package chat

import (
	"context"
	"log"

	"github.com/WithSoull/ChatServer/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) DeleteChat(ctx context.Context, senderID, chatID int64) error {
	_, err := s.chatRepo.Get(ctx, chatID)
	if err != nil {
		log.Printf("[Service Layer] Failed to get chat %d: %v", chatID, err)
		return status.Errorf(codes.NotFound, "chat does not exist")
	}

	ok, role := s.chatParticipantRepo.GetUserRole(ctx, chatID, senderID)
	if ok && role == model.ROLE_OWNER {
		if err := s.chatRepo.Delete(ctx, chatID); err != nil {
			log.Printf("[Service Layer] Failed to delete chat %d: %v", chatID, err)
			return status.Errorf(codes.Internal, "failed to delete chat")
		}
		return nil
	}

	return status.Errorf(codes.PermissionDenied, "only owner of the chat can delete it")
}
