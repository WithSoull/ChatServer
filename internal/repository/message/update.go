package message

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/platform_common/pkg/client/db"
	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/ChatServer/internal/repository/converter"
)

func (r *messageRepo) Update(ctx context.Context, msg model.Message) error {
	rmsg := converter.FromModelToRepoMessage(msg)

	builder := sq.Update(TableName).
		PlaceholderFormat(sq.Dollar).
		Set(TextColumn, rmsg.Text).
		Set(IsPinnedColumn, rmsg.IsPinned).
		Set(UpdatedAtColumn, time.Now()).
		Where(sq.Eq{IDColumn: rmsg.ID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "message_repository:Update",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return domainerrors.ErrMessageNotFound
	}

	return nil
}
