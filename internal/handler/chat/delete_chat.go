package chat

import (
	"context"

	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) DeleteChat(ctx context.Context, req *desc.DeleteChatRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, h.service.DeleteChat(ctx, req.GetChatId())
}
