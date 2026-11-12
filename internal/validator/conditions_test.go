package validator

import (
	"context"
	"errors"
	"strings"
	"testing"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
)

type conditionCase struct {
	name    string
	input   string
	wantErr error
}

func runConditionCases(t *testing.T, newCond func(string) validate.Condition, cases []conditionCase) {
	t.Helper()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := newCond(tc.input)(context.Background())
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestValidateNotEmptyMessage(t *testing.T) {
	runConditionCases(t, ValidateNotEmptyMessage, []conditionCase{
		{name: "empty", input: "", wantErr: domainerrors.ErrEmptyMessageText},
		{name: "single space", input: " ", wantErr: domainerrors.ErrEmptyMessageText},
		{name: "many spaces", input: "     ", wantErr: domainerrors.ErrEmptyMessageText},
		{name: "tabs and newlines", input: "\t\n\r\n ", wantErr: domainerrors.ErrEmptyMessageText},
		{name: "unicode whitespace", input: "  　", wantErr: domainerrors.ErrEmptyMessageText},
		{name: "plain text", input: "hello"},
		{name: "text surrounded by spaces", input: "  hi  "},
		{name: "single char", input: "a"},
		{name: "cyrillic", input: "привет"},
	})
}

func TestValidateMessageMaxLength(t *testing.T) {
	runConditionCases(t, ValidateMessageMaxLength, []conditionCase{
		{name: "empty", input: ""},
		{name: "short", input: "hello"},
		{name: "exactly max ascii", input: strings.Repeat("a", MaxMessageLength)},
		{name: "max plus one ascii", input: strings.Repeat("a", MaxMessageLength+1), wantErr: domainerrors.ErrMessageTooLong},
		// Multi-byte runes: 4096 runes is more than 4096 bytes but must still pass.
		{name: "exactly max cyrillic", input: strings.Repeat("я", MaxMessageLength)},
		{name: "exactly max emoji", input: strings.Repeat("😀", MaxMessageLength)},
		{name: "max plus one cyrillic", input: strings.Repeat("я", MaxMessageLength+1), wantErr: domainerrors.ErrMessageTooLong},
		{name: "way too long", input: strings.Repeat("a", MaxMessageLength*3), wantErr: domainerrors.ErrMessageTooLong},
	})
}

func TestValidateEmptyName(t *testing.T) {
	runConditionCases(t, ValidateEmptyName, []conditionCase{
		{name: "empty", input: "", wantErr: domainerrors.ErrEmptyChatName},
		{name: "single space", input: " ", wantErr: domainerrors.ErrEmptyChatName},
		{name: "many spaces", input: "    ", wantErr: domainerrors.ErrEmptyChatName},
		{name: "tabs and newlines", input: "\t\n", wantErr: domainerrors.ErrEmptyChatName},
		{name: "plain name", input: "general"},
		{name: "name with spaces inside", input: "team chat"},
		{name: "name surrounded by spaces", input: " team "},
	})
}

func TestValidateNoChangesProvided(t *testing.T) {
	runConditionCases(t, ValidateNoChangesProvided, []conditionCase{
		{name: "empty", input: "", wantErr: domainerrors.ErrNoChangesProvided},
		{name: "only spaces", input: "   ", wantErr: domainerrors.ErrNoChangesProvided},
		{name: "new text", input: "edited"},
	})
}

func TestValidateCombined(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{name: "valid", text: "hello", wantErr: false},
		{name: "whitespace only", text: " ", wantErr: true},
		{name: "too long", text: strings.Repeat("a", MaxMessageLength+1), wantErr: true},
		{name: "too long whitespace", text: strings.Repeat(" ", MaxMessageLength+1), wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Validate(context.Background(),
				ValidateNotEmptyMessage(tc.text),
				ValidateMessageMaxLength(tc.text),
			)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v, got %v", tc.wantErr, err)
			}
			if err != nil && !validate.IsValidationError(err) {
				t.Fatalf("expected validation error, got %T: %v", err, err)
			}
		})
	}
}
