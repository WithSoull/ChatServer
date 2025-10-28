package chat

import (
	"context"

	"github.com/WithSoull/ChatServer/internal/model"
)

func (s *Service) ConnectChat(ctx context.Context, senderID, chatID int64) (chan *model.Message, error) {
	if err := s.userExist(ctx, senderID); err != nil {
		return nil, err
	}

	_, err := s.checkUserRole(ctx, chatID, senderID, model.ROLE_USER)
	if err != nil {
		return nil, err
	}

	return s.streams.MsgStream(chatID, senderID), nil
}

// WARN: Dont check user permisions !!!
func (s *Service) DisconnectChat(ctx context.Context, senderID, chatID int64) {
	s.streams.RemoveMsgStream(chatID, senderID)
}
