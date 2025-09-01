package chat

import (
	"context"

	converter "github.com/WithSoull/ChatServer/internal/converter"
	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) UpdateUserRole(ctx context.Context, req *desc.UpdateUserRoleRequest) (*emptypb.Empty, error) {
	err := h.service.UpdateUserRole(ctx, req.GetSenderId(), req.GetChatId(), req.GetUserId(), converter.FromProtoToModelRole(req.Role))

	return &emptypb.Empty{}, err
}
