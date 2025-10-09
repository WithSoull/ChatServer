package validator

import (
	"context"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
)

func ValidateNotEmptyMessege(text string) validate.Condition {
	return func(ctx context.Context) error {
		if text == "" {
			return domainerrors.ErrEmptyMessageText
		}

		return nil
	}
}

func ValidateNoChangesProvided(text string) validate.Condition {
	return func(ctx context.Context) error {
		if text == "" {
			return domainerrors.ErrNoChangesProvided
		}

		return nil
	}
}
