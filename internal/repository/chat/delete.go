package chat

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/platform_common/pkg/client/db"
	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
)

func (r *chatRepo) Delete(ctx context.Context, chatID int64) error {
	builder := sq.Delete(TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": chatID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_repository:Delete",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return domainerrors.ErrChatNotFound
	}

	return nil
}
