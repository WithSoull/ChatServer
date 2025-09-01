package chat

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetChat(ctx context.Context, senderID, chatID int64) (model.Chat, []model.Message, error) {
	if ok, _ := s.chatParticipantRepo.GetUserRole(ctx, chatID, senderID); !ok {
		return model.Chat{}, nil, status.Errorf(codes.PermissionDenied, "you are not a member of this chat")
	}

	chat, err := s.chatRepo.Get(ctx, chatID)
	if err != nil {
		return model.Chat{}, nil, status.Errorf(codes.NotFound, "chat not found")
	}

	chatParticipants, err := s.chatParticipantRepo.GetUsers(ctx, chatID)
	if err != nil {
		return model.Chat{}, nil, status.Errorf(codes.Internal, "failed to get users of the chat")
	}
	chat.ChatInfo.UserIDs = chatParticipants

	msgs, err := s.msgRepo.GetByChat(ctx, chatID)
	if err != nil {
		return model.Chat{}, nil, status.Errorf(codes.Internal, "failed to get messages of the chat")
	}

	return chat, msgs, nil
}
