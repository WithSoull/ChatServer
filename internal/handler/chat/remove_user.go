package chat

import (
	"context"

	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) RemoveUser(ctx context.Context, req desc.RemoveUserRequest) (*emptypb.Empty, error) {
	err := h.service.RemoveUser(ctx, req.GetSenderId(), req.GetChatId(), req.GetUserId())

	return &emptypb.Empty{}, err
}
