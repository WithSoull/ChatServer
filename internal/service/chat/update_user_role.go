package chat

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) UpdateUserRole(ctx context.Context, senderID, chatID, userID int64, newRole model.Role) error {
	senderRole, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_ADMIN)
	if err != nil {
		return err
	}

	targetRole, err := s.checkUserRole(ctx, chatID, userID, model.ROLE_USER)
	if err != nil {
		return err
	}

	if newRole >= senderRole {
		return status.Errorf(codes.PermissionDenied, "cannot assign a role higher or equal than your own")
	}

	if targetRole >= senderRole {
		log.Printf("[DEBUG] targetRole = %d | senderRole = %d ", targetRole, senderRole)
		return status.Errorf(codes.PermissionDenied, "cannot change role of a user with equal or higher role")
	}

	if senderRole == model.ROLE_OWNER && newRole == model.ROLE_OWNER {
		return status.Errorf(codes.InvalidArgument, "cannot assign owner role to another member")
	}

	if err := s.chatParticipantRepo.UpdateUserRole(ctx, chatID, userID, newRole); err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			return status.Errorf(codes.Internal, "failed to update user role")
		}
		return grpcErr
	}

	return nil
}

func (s *Service) checkUserRole(ctx context.Context, chatID, userID int64, neededRole model.Role) (model.Role, error) {
	senderRole, err := s.chatParticipantRepo.GetUserRole(ctx, chatID, userID)
	if err != nil {
		if err == domainerrors.ErrCantDefineWhoDoesNotExistUserOrChat {
			if _, err := s.chatRepo.Get(ctx, chatID); err != nil {
				isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
				if isLogNeeded {
					log.Printf("[Service Layer] failed to get chat: %v", err)
				}
				return model.ROLE_USER, grpcErr
			}
			return model.ROLE_USER, status.Errorf(codes.PermissionDenied, "user(ID=%d) are not a member of the chat", userID)
		}

		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to get user(ID=%d) role from chat(ID=%d): %v", userID, chatID, err)
		}
		return model.ROLE_USER, grpcErr
	}

	if senderRole < neededRole {
		switch neededRole {
		case model.ROLE_USER:
			log.Print("[Service Layer] send sender role < needed role")
			return model.ROLE_USER, status.Error(codes.Internal, "failed to check users's role")
		case model.ROLE_ADMIN:
			return model.ROLE_USER, status.Error(codes.PermissionDenied, "only admins and owner has permission to make this action")
		case model.ROLE_OWNER:
			return model.ROLE_USER, status.Error(codes.PermissionDenied, "only owner has permission to make this action")
		default:
			return model.ROLE_USER, status.Errorf(codes.PermissionDenied, "user(ID=%d) have no permission do it", userID)
		}
	}

	return senderRole, nil
}
