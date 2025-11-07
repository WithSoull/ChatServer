package chat

import (
	"context"

	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) DeleteMessage(ctx context.Context, req *desc.DeleteMessageRequest) (*emptypb.Empty, error) {
	err := h.service.DeleteMessage(ctx, req.GetMessageId())

	return &emptypb.Empty{}, err
}
