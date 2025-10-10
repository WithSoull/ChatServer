package message

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/platform_common/pkg/client/db"
	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/ChatServer/internal/repository/converter"
	"github.com/jackc/pgx/v4"
)

func (r *messageRepo) Create(ctx context.Context, msg model.Message) (int64, error) {
	rmsg := converter.FromModelToRepoMessage(msg)
	now := time.Now()

	builder := sq.Insert(TableName).
		PlaceholderFormat(sq.Dollar).
		Columns(ChatIDColumn, SenderIDColumn, TextColumn, IsPinnedColumn, SendAtColumn, UpdatedAtColumn).
		Values(rmsg.ChatID, rmsg.SenderID, rmsg.Text, rmsg.IsPinned, now, now).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "message_repository:Create",
		QueryRaw: query,
	}

	var msgID int64

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&msgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domainerrors.ErrMessageNotFound
		}
		return 0, err
	}

	return msgID, nil
}
