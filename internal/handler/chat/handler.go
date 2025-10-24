package chat

import (
	"github.com/WithSoull/ChatServer/internal/service"
	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
)

type Chat struct {
}

type Handler struct {
	desc.UnimplementedChatV1Server
	service service.ChatService
}

func NewHandler(service service.ChatService) *Handler {
	return &Handler{
		service: service,
	}
}
