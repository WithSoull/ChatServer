package chat

import (
	"context"
	"log"

	"github.com/WithSoull/ChatServer/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) UpdateUserRole(ctx context.Context, senderID, chatID, userID int64, newRole model.Role) error {
	ok, senderRole := s.chatParticipantRepo.GetUserRole(ctx, chatID, senderID)
	if !ok {
		return status.Errorf(codes.PermissionDenied, "you are not a member of the chat")
	}

	ok, targetRole := s.chatParticipantRepo.GetUserRole(ctx, chatID, userID)
	if !ok {
		return status.Errorf(codes.NotFound, "target user is not a member of the chat")
	}

	if newRole > senderRole {
		return status.Errorf(codes.PermissionDenied, "cannot assign a role higher than your own")
	}

	if targetRole >= senderRole {
		return status.Errorf(codes.PermissionDenied, "cannot change role of a user with equal or higher role")
	}

	if senderRole == model.ROLE_OWNER && newRole == model.ROLE_OWNER {
		return status.Errorf(codes.InvalidArgument, "cannot assign owner role to another member")
	}

	if err := s.chatParticipantRepo.UpdateUserRole(ctx, chatID, userID, newRole); err != nil {
		log.Printf("failed to update role of user %d in chat %d: %v", userID, chatID, err)
		return status.Errorf(codes.Internal, "failed to update user role")
	}

	return nil
}
