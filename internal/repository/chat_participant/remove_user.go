package chat_participant

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/platform_common/pkg/client/db"
)

func (r *chatParticipantRepo) RemoveUserFromChat(ctx context.Context, chatID, userID int64) error {
	builder := sq.Delete(TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{ChatIDColumn: chatID, UserIDColumn: userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_participant_repository:RemoveUser",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return domainerrors.ErrUserNotFound
	}

	return nil
}

func (r *chatParticipantRepo) RemoveUserFromAllChats(ctx context.Context, userID int64) error {
	builder := sq.Delete(TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{UserIDColumn: userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_participant_repository:RemoveUserFromAllChats",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return domainerrors.ErrUserNotFound
	}
	return nil
}
