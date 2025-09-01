package chat

import (
	"context"
	"log"

	"github.com/WithSoull/ChatServer/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) RemoveUser(ctx context.Context, senderID, chatID, userID int64) error {
	ok, senderRole := s.chatParticipantRepo.GetUserRole(ctx, chatID, senderID)
	if !ok {
		return status.Errorf(codes.PermissionDenied, "you are not a member of the chat")
	}
	if senderRole < 1 {
		return status.Errorf(codes.PermissionDenied, "only admins and owner of the chat can remove users")
	}

	if ok, targetRole := s.chatParticipantRepo.GetUserRole(ctx, chatID, userID); ok && targetRole == model.ROLE_OWNER {
		return status.Errorf(codes.PermissionDenied, "owner of the chat cannot be removed")
	}

	if err := s.chatParticipantRepo.RemoveUser(ctx, chatID, userID); err != nil {
		log.Printf("[Service Layer] failed to remove user %d from chat %d: %v", userID, chatID, err)
		return status.Errorf(codes.Internal, "failed to remove user from chat")
	}

	return nil
}
