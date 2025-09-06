package chat

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

func (r *chatRepo) Get(ctx context.Context, chatID int64) (model.Chat, error) {
	builder := sq.Select(IDColumn, OwnerIDColumn, NameColumn, DescriptionColumn, CreatedAtColumn, UpdatedAtColumn).
		From(TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": chatID})

	query, args, err := builder.ToSql()
	if err != nil {
		return model.Chat{}, err
	}

	q := db.Query{
		Name:     "chat_repository:Get",
		QueryRaw: query,
	}

	var chat rmodel.Chat

	err = r.db.DB().ScanOneContext(ctx, &chat, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Chat{}, domainerrors.ErrChatNotFound
		}
		return model.Chat{}, err
	}

	return converter.FromRepoToModelChat(chat, nil), nil
}
