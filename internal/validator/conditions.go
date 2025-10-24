package validator

import (
	"context"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
)

func ValidateNotEmptyMessage(text string) validate.Condition {
	return func(ctx context.Context) error {
		if text == "" {
			return domainerrors.ErrEmptyMessageText
		}

		return nil
	}
}

func ValidateEmptyName(name string) validate.Condition {
	return func(ctx context.Context) error {
		if name == "" {
			return domainerrors.ErrEmptyChatName
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
