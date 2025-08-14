package chat

import (
	"github.com/WithSoull/ChatServer/internal/service"
	desc "github.com/WithSoull/ChatServer/pkg/chat/v1"
)

type Handler struct {
	desc.UnimplementedChatV1Server
	service service.ChatService
}
