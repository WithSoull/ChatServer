package message

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/ChatServer/internal/client/db"
	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/ChatServer/internal/repository/converter"
	rmodel "github.com/WithSoull/ChatServer/internal/repository/model"
	"github.com/jackc/pgx/v4"
)

func (r *messageRepo) Get(ctx context.Context, messageID int64) (model.Message, error) {
	builder := sq.Select(IDColumn, ChatIDColumn, SenderIDColumn, TextColumn, IsPinnedColumn, SendAtColumn, UpdatedAtColumn).
		From(TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{IDColumn: messageID})

	query, args, err := builder.ToSql()
	if err != nil {
		return model.Message{}, err
	}

	q := db.Query{
		Name:     "message_repository:Get",
		QueryRaw: query,
	}

	var rmsg rmodel.Message
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&rmsg.ID,
		&rmsg.ChatID,
		&rmsg.SenderID,
		&rmsg.Text,
		&rmsg.IsPinned,
		&rmsg.SendAt,
		&rmsg.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Message{}, domainerrors.ErrMessageNotFound
		}
		return model.Message{}, err
	}

	return converter.FromRepoToModelMessage(rmsg), nil
}
