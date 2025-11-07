package chat

import (
	"context"

	converter "github.com/WithSoull/ChatServer/internal/converter"
	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
)

func (h *Handler) GetChat(ctx context.Context, req *desc.GetChatRequest) (*desc.GetChatResponse, error) {
	chat, msgs, err := h.service.GetChat(ctx, req.GetChatId())
	if err != nil {
		return nil, err
	}

	return &desc.GetChatResponse{
		Chat:     converter.FromModelToProtoChat(chat),
		Messages: converter.FromModelToProtoMessages(msgs),
	}, nil
}
