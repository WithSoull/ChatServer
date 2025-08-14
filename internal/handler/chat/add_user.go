package chat

import (
	"context"

	conventer "github.com/WithSoull/ChatServer/internal/conventer/chat"
	chat_v1 "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) AddUser(ctx context.Context, req *chat_v1.AddUserRequest) (*emptypb.Empty, error) {
	err := h.service.AddUser(ctx, req.GetSenderId(), req.GetChatId(), req.GetUserId(), conventer.FromProtoToModelRole(req.GetRole()))

	return &emptypb.Empty{}, err
}
