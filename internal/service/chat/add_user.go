package chat

import (
	"context"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
)

func (s *Service) AddUser(ctx context.Context, chatID, userID int64, role model.Role) error {
	senderID, ok := claimsctx.ExtractUserID(ctx)
	if !ok {
		return domainerrors.ErrFailedToVerify
	}

	if err := s.userExist(ctx, userID); err != nil {
		return err
	}

	// role > 1 is equal admin or owner
	senderRole, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_ADMIN)
	if err != nil {
		return err
	}
	if role >= senderRole {
		return domainerrors.ErrCannotAssignHigherRole
	}

	err = s.chatParticipantRepo.AddUser(ctx, chatID, userID, role)
	if err != nil {
		return err
	}

	return nil
}
