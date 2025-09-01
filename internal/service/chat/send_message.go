package chat

import (
	"context"
	"log"

	"github.com/WithSoull/ChatServer/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) SendMessage(ctx context.Context, senderID int64, chatID int64, text string) error {
	ok, _ := s.chatParticipantRepo.GetUserRole(ctx, chatID, senderID)
	if !ok {
		return status.Errorf(codes.PermissionDenied, "you must be a member of the chat to send a message")
	}

	if text == "" {
		return status.Errorf(codes.InvalidArgument, "message text cannot be empty")
	}

	msg := model.Message{
		SenderID: senderID,
		ChatID:   chatID,
		Text:     text,
	}

	_, err := s.msgRepo.Create(ctx, msg)
	if err != nil {
		log.Printf("failed to create message in chat %d from sender %d: %v", chatID, senderID, err)
		return status.Errorf(codes.Internal, "failed to send message")
	}

	return nil
}
