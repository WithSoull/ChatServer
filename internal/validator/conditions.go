package validator

import (
	"context"
	"strings"
	"unicode/utf8"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
)

const (
	// MaxMessageLength is the maximum allowed message length in runes.
	MaxMessageLength = 4096
	// MaxChatNameLength is the maximum allowed chat name length in runes.
	MaxChatNameLength = 128
)

func ValidateNotEmptyMessage(text string) validate.Condition {
	return func(ctx context.Context) error {
		if strings.TrimSpace(text) == "" {
			return domainerrors.ErrEmptyMessageText
		}

		return nil
	}
}

func ValidateMessageMaxLength(text string) validate.Condition {
	return func(ctx context.Context) error {
		if utf8.RuneCountInString(text) > MaxMessageLength {
			return domainerrors.ErrMessageTooLong
		}

		return nil
	}
}

func ValidateEmptyName(name string) validate.Condition {
	return func(ctx context.Context) error {
		if strings.TrimSpace(name) == "" {
			return domainerrors.ErrEmptyChatName
		}

		return nil
	}
}

func ValidateChatNameMaxLength(name string) validate.Condition {
	return func(ctx context.Context) error {
		if utf8.RuneCountInString(name) > MaxChatNameLength {
			return domainerrors.ErrChatNameTooLong
		}

		return nil
	}
}

func ValidateNoChangesProvided(text string) validate.Condition {
	return func(ctx context.Context) error {
		if strings.TrimSpace(text) == "" {
			return domainerrors.ErrNoChangesProvided
		}

		return nil
	}
}
