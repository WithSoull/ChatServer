package chat

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) EditMessage(
	ctx context.Context,
	senderID, messageID int64,
	newText *string,
	newIsPinned *bool,
) error {
	msg, err := s.msgRepo.Get(ctx, messageID)
	if err != nil {
		log.Printf("failed to get message %d: %v", messageID, err)
		return status.Errorf(codes.NotFound, "message not found")
	}

	ok, senderRole := s.chatParticipantRepo.GetUserRole(ctx, msg.ChatID, senderID)
	if !ok {
		return status.Errorf(codes.PermissionDenied, "you are not a member of the chat")
	}

	if newText == nil && newIsPinned == nil {
		return status.Errorf(codes.InvalidArgument, "no changes provided")
	}

	if newText != nil {
		if senderID != msg.SenderID {
			return status.Errorf(codes.PermissionDenied, "you can edit only your own messages")
		}
		msg.Text = *newText
	}

	if newIsPinned != nil {
		if senderRole < 1 {
			return status.Errorf(codes.PermissionDenied, "only admins or owner can pin/unpin messages")
		}
		msg.IsPinned = *newIsPinned
	}

	if err := s.msgRepo.Update(ctx, msg); err != nil {
		log.Printf("failed to update message %d: %v", messageID, err)
		return status.Errorf(codes.Internal, "failed to update message")
	}

	return nil
}
