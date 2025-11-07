package chat

import (
	"github.com/WithSoull/ChatServer/internal/service"
	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
	"github.com/WithSoull/platform_common/pkg/tokens"
)

type Chat struct {
}

type Handler struct {
	desc.UnimplementedChatV1Server
	service       service.ChatService
	tokenVerifier tokens.TokenVerifier
}

func NewHandler(service service.ChatService, tokenVerifier tokens.TokenVerifier) *Handler {
	return &Handler{
		service:       service,
		tokenVerifier: tokenVerifier,
	}
}
