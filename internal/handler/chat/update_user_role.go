package chat

import (
	"context"

	conventer "github.com/WithSoull/ChatServer/internal/conventer/chat"
	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) UpdateUserRole(ctx context.Context, req *desc.UpdateUserRoleRequest) (*emptypb.Empty, error) {
	err := h.service.UpdateUserRole(ctx, req.GetSenderId(), req.GetChatId(), req.GetUserId(), conventer.FromProtoToModelRole(req.Role))

	return &emptypb.Empty{}, err
}
