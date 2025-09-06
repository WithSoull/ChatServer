package chat

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) RemoveUser(ctx context.Context, senderID, chatID, userID int64) error {
	targetRole, err := s.checkUserRole(ctx, chatID, userID, model.ROLE_USER)
	if err != nil {
		return err
	}

	senderRole, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_ADMIN)
	if err != nil {
		return err
	}

	if senderRole <= targetRole {
		return status.Error(codes.PermissionDenied, "cannot remove a user with equal or higher role")
	}

	if err := s.chatParticipantRepo.RemoveUser(ctx, chatID, userID); err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to remove user %d from chat %d: %v", userID, chatID, err)
		}
		return grpcErr
	}

	return nil
}
