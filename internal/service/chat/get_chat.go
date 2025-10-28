package chat

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/platform_common/pkg/tracing"
)

func (s *Service) GetChat(ctx context.Context, senderID, chatID int64) (model.Chat, []model.Message, error) {
	if err := s.userExist(ctx, senderID); err != nil {
		return model.Chat{}, nil, err
	}

	if _, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_USER); err != nil {
		return model.Chat{}, nil, err
	}

	ctx, span := tracing.StartSpan(ctx, "repo:chats:Get")
	chat, err := s.chatRepo.Get(ctx, chatID)
	if err != nil {
		return model.Chat{}, nil, err
	}
	span.End()

	ctx, span = tracing.StartSpan(ctx, "repo:chats:GetUsers")
	chatParticipants, err := s.chatParticipantRepo.GetUsers(ctx, chatID)
	if err != nil {
		return model.Chat{}, nil, err
	}
	chat.ChatInfo.UserIDs = chatParticipants
	span.End()

	ctx, span = tracing.StartSpan(ctx, "repo:messages:GetByChat")
	msgs, err := s.msgRepo.GetByChat(ctx, chatID)
	if err != nil {
		return model.Chat{}, nil, err
	}
	span.End()

	return chat, msgs, nil
}
