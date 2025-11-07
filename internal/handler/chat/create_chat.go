package chat

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/converter"
	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
)

func (h *Handler) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	chatId, err := h.service.CreateChat(ctx, converter.FromProtoToModelChatInfo(req.ChatInfo))

	return &desc.CreateChatResponse{
		ChatId: chatId,
	}, err
}
