package chat

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) AddUser(ctx context.Context, senderID, chatID, userID int64, role model.Role) error {
	// role > 1 is equal admin or owner
	senderRole, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_ADMIN)
	if err != nil {
		return err
	}
	if role >= senderRole {
		return status.Errorf(codes.PermissionDenied, "cannot assign a role higher or equal than your own")
	}

	err = s.chatParticipantRepo.AddUser(ctx, chatID, userID, role)
	if err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to add user %d to chat %d: %v", userID, chatID, err)
		}
		return grpcErr
	}

	return nil
}
