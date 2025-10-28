package chat

import (
	"github.com/WithSoull/ChatServer/internal/converter"
	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
)

func (h *Handler) ConnectChat(req *desc.ConnectChatRequest, stream desc.ChatV1_ConnectChatServer) error {
	_, old_messages, err := h.service.GetChat(stream.Context(), req.GetSenderId(), req.GetChatId())
	if err != nil {
		return err
	}

	for _, msg := range converter.FromModelToProtoMessages(old_messages) {
		if err := stream.Send(msg); err != nil {
			return err
		}
	}

	channel, err := h.service.ConnectChat(stream.Context(), req.GetSenderId(), req.GetChatId())
	if err != nil {
		return nil
	}

	for {
		select {
		case msg, ok := <-channel:
			if !ok || msg == nil {
				logger.Debug(stream.Context(), "connection was closed because chat was deleted", zap.Int64("chatID", req.GetChatId()), zap.Int64("senderID", req.GetSenderId()))
				return domainerrors.ErrChatHasBeenDeleted
			}

			if err := stream.Send(converter.FromModelToProtoMessage(*msg)); err != nil {
				return err
			}
		case <-stream.Context().Done():
			h.service.DisconnectChat(stream.Context(), req.GetSenderId(), req.GetChatId())
			logger.Debug(stream.Context(), "connection was closed", zap.Int64("chatID", req.GetChatId()), zap.Int64("senderID", req.GetSenderId()))
			return nil
		}
	}
}
