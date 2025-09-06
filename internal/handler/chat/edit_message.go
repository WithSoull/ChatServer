package chat

import (
	"context"

	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) EditMessage(ctx context.Context, req *desc.EditMessageRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, h.service.EditMessage(ctx, req.GetSenderId(), req.GetMessageId(), req.GetNewText())

}
