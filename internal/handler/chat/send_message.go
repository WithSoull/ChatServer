package chat

import (
	"context"

	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, h.service.SendMessage(ctx, req.GetChatId(), req.GetText())
}
