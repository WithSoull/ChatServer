package validator

import (
	"context"
	"errors"
	"strings"
	"testing"

	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
)

func TestValidateNotEmptyMessage(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr error
	}{
		{name: "empty", text: "", wantErr: domainerrors.ErrEmptyMessageText},
		{name: "spaces only", text: "   ", wantErr: domainerrors.ErrEmptyMessageText},
		{name: "mixed whitespace", text: " \t\n\r ", wantErr: domainerrors.ErrEmptyMessageText},
		{name: "unicode whitespace", text: "  ", wantErr: domainerrors.ErrEmptyMessageText},
		{name: "plain text", text: "hello", wantErr: nil},
		{name: "text with surrounding spaces", text: "  hi  ", wantErr: nil},
		{name: "cyrillic", text: "привет", wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNotEmptyMessage(tt.text)(context.Background())
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ValidateNotEmptyMessage(%q) = %v, want %v", tt.text, err, tt.wantErr)
			}
		})
	}
}

func TestValidateMessageLength(t *testing.T) {
	limit := domainerrors.MaxMessageLength

	tests := []struct {
		name    string
		text    string
		wantErr error
	}{
		{name: "empty", text: "", wantErr: nil},
		{name: "ascii at limit", text: strings.Repeat("a", limit), wantErr: nil},
		{name: "ascii over limit", text: strings.Repeat("a", limit+1), wantErr: domainerrors.ErrMessageTooLong},
		// 4096 Cyrillic runes take 8192 bytes and must still pass.
		{name: "cyrillic at limit", text: strings.Repeat("я", limit), wantErr: nil},
		{name: "cyrillic over limit", text: strings.Repeat("я", limit+1), wantErr: domainerrors.ErrMessageTooLong},
		{name: "emoji at limit", text: strings.Repeat("😀", limit), wantErr: nil},
		{name: "emoji over limit", text: strings.Repeat("😀", limit+1), wantErr: domainerrors.ErrMessageTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMessageLength(tt.text)(context.Background())
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ValidateMessageLength(%d bytes) = %v, want %v", len(tt.text), err, tt.wantErr)
			}
		})
	}
}
