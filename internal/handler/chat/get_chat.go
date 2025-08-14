package chat

import (
	"context"

	conventer "github.com/WithSoull/ChatServer/internal/conventer/chat"
	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
)

func (h *Handler) GetChat(ctx context.Context, req *desc.GetChatRequest) (*desc.GetChatResponse, error) {
	chat, msgs, err := h.service.GetChat(ctx, req.GetSenderId(), req.GetChatId())
	if err != nil {
		return nil, err
	}

	return &desc.GetChatResponse{
		Chat:     conventer.FromModelToProtoChat(chat),
		Messages: conventer.FromModelToProtoMessages(msgs),
	}, nil
}
