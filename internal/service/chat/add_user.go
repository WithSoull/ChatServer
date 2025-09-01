package chat

import (
	"context"
	"errors"
	"log"

	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/jackc/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) AddUser(ctx context.Context, senderID, chatID, userID int64, role model.Role) error {
	// role > 1 is equal admin or owner
	ok, senderRole := s.chatParticipantRepo.GetUserRole(ctx, chatID, senderID)
	if !ok {
		return status.Errorf(codes.PermissionDenied, "you are not a member of the chat")
	}

	if senderRole < 1 {
		return status.Errorf(codes.PermissionDenied, "only admins and owner of the chat can add new users")
	}

	if role == model.ROLE_OWNER {
		return status.Errorf(codes.InvalidArgument, "only 1 owner can be in the chat")
	}

	err := s.chatParticipantRepo.AddUser(ctx, chatID, userID, role)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return status.Errorf(codes.AlreadyExists, "user already in the chat, for updating role use update_role endpoint")
		}

		log.Printf("[Service Layer] failed to add user %d to chat %d: %v", userID, chatID, err)
		return status.Errorf(codes.Internal, "failed to add user to chat")
	}

	return nil
}
