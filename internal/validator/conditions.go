package validator

import (
	"context"
	"strings"
	"unicode/utf8"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
)

func ValidateNotEmptyMessage(text string) validate.Condition {
	return func(ctx context.Context) error {
		if strings.TrimSpace(text) == "" {
			return domainerrors.ErrEmptyMessageText
		}

		return nil
	}
}

// ValidateMessageLength counts runes, not bytes, so Cyrillic and other
// multi-byte text gets the same limit as ASCII.
func ValidateMessageLength(text string) validate.Condition {
	return func(ctx context.Context) error {
		if utf8.RuneCountInString(text) > domainerrors.MaxMessageLength {
			return domainerrors.ErrMessageTooLong
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
