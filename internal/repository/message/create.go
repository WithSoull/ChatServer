package message

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/ChatServer/internal/repository/converter"
	"github.com/WithSoull/platform_common/pkg/client/db"
)

func (r *messageRepo) Create(ctx context.Context, msg *model.Message) error {
	rmsg := converter.FromModelToRepoMessage(*msg)
	now := time.Now()

	builder := sq.Insert(TableName).
		PlaceholderFormat(sq.Dollar).
		Columns(ChatIDColumn, SenderIDColumn, TextColumn, IsPinnedColumn, SendAtColumn, UpdatedAtColumn).
		Values(rmsg.ChatID, rmsg.SenderID, rmsg.Text, rmsg.IsPinned, now, now).
		Suffix(fmt.Sprintf("RETURNING %s, %s, %s", IDColumn, SendAtColumn, UpdatedAtColumn))

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "message_repository:Create",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&rmsg.ID, &rmsg.SendAt, &rmsg.UpdatedAt)
	if err != nil {
		return err
	}

	msg.SendAt = rmsg.SendAt
	msg.MessageID = rmsg.ID
	msg.UpdatedAt = rmsg.UpdatedAt

	return nil
}
