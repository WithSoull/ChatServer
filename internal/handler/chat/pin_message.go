package chat

import (
	"context"

	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) PinMessage(ctx context.Context, req *desc.PinMessageRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, h.service.PinMessage(ctx, req.GetMessageId(), req.GetNewIsPinned())

}
