package chat

import (
	"context"

	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
)

func (h *Handler) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	chatId, err := h.service.CreateChat(
		ctx,
		req.GetSenderId(),
		req.GetChatInfo().GetName(),
		req.GetChatInfo().GetDescription(),
		req.GetChatInfo().GetUserIds(),
	)

	return &desc.CreateChatResponse{
		ChatId: chatId,
	}, err
}
