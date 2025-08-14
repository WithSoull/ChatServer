package chat

import (
	"github.com/WithSoull/ChatServer/internal/client/db"
	"github.com/WithSoull/ChatServer/internal/repository"
)

type Service struct {
	chatRepo repository.ChatRepo
	msgRepo  repository.MessageRepo

	txManager db.TxManager
}
