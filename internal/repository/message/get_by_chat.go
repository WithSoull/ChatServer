package message

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/platform_common/pkg/client/db"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/ChatServer/internal/repository/converter"
	rmodel "github.com/WithSoull/ChatServer/internal/repository/model"
)

func (r *messageRepo) GetByChat(ctx context.Context, chatID int64) ([]model.Message, error) {
	builder := sq.Select(IDColumn, ChatIDColumn, SenderIDColumn, TextColumn, IsPinnedColumn, SendAtColumn, UpdatedAtColumn).
		From(TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{ChatIDColumn: chatID}).
		OrderBy(SendAtColumn + " ASC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "message_repository:GetByChat",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rmessages []rmodel.Message
	for rows.Next() {
		var rmsg rmodel.Message
		if err := rows.Scan(
			&rmsg.ID,
			&rmsg.ChatID,
			&rmsg.SenderID,
			&rmsg.Text,
			&rmsg.IsPinned,
			&rmsg.SendAt,
			&rmsg.UpdatedAt,
		); err != nil {
			return nil, err
		}
		rmessages = append(rmessages, rmsg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return converter.FromRepoToModelMessages(rmessages), nil
}
