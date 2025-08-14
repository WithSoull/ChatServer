package chat

import "context"

func (s *Service) Delete(ctx context.Context, chatID int64) error {
	return s.repo.DeleteChat(ctx, chatID)
}

