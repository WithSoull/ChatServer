package chat

import (
	"context"
	"errors"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/platform_common/pkg/logger"
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
		return domainerrors.ErrCannotAssignHigherRole
	}

	if targetRole >= senderRole {
		return domainerrors.ErrCannotAssignHigherRole
	}

	if senderRole == model.ROLE_OWNER && newRole == model.ROLE_OWNER {
		return domainerrors.ErrCannotAssignOwnerRole
	}

	if err := s.chatParticipantRepo.UpdateUserRole(ctx, chatID, userID, newRole); err != nil {
		return err
	}

	return nil
}

func (s *Service) checkUserRole(ctx context.Context, chatID, userID int64, neededRole model.Role) (model.Role, error) {
	senderRole, err := s.chatParticipantRepo.GetUserRole(ctx, chatID, userID)
	if err != nil {
		if errors.Is(err, domainerrors.ErrCantDefineWhoDoesNotExistUserOrChat) {
			if _, err := s.chatRepo.Get(ctx, chatID); err != nil {
				return model.ROLE_USER, err
			}
			return model.ROLE_USER, domainerrors.ErrUserNotMember
		}

		return model.ROLE_USER, err
	}

	if senderRole < neededRole {
		switch neededRole {
		case model.ROLE_USER:
			logger.Error(ctx, "send sender role < needed role")
			return model.ROLE_USER, domainerrors.ErrFailedToCheckRole
		case model.ROLE_ADMIN:
			return model.ROLE_USER, domainerrors.ErrOnlyAdminsAllowed
		case model.ROLE_OWNER:
			return model.ROLE_USER, domainerrors.ErrOnlyOwnerAllowed
		default:
			return model.ROLE_USER, domainerrors.ErrUserNoPermission
		}
	}

	return senderRole, nil
}
