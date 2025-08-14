package chat

import "context"

func (s *Service) Create(ctx context.Context, userIDs []int64) (int64, error) {

	chat_id, err := s.repo.CreateChat(ctx, userIDs)

	// NEW FEATURE: we can add names for the chat and censor obscene language

	return chat_id, err
}

