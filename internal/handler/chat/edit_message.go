package chat

import (
	"context"

	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) EditMessage(ctx context.Context, req *desc.EditMessageRequest) (*emptypb.Empty, error) {
	var newText *string
	if req.GetNewText() != nil {
		newText = &req.GetNewText().Value
	}

	var newIsPinned *bool
	if req.GetNewIsPinned() != nil {
		newIsPinned = &req.GetNewIsPinned().Value
	}

	err := h.service.EditMessage(ctx, req.GetSenderId(), req.GetMessageId(), newText, newIsPinned)

	return &emptypb.Empty{}, err
}
