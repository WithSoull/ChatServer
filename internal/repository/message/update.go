package message

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/ChatServer/internal/client/db"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/ChatServer/internal/repository/converter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		log.Printf("failed to update message %d: %v", rmsg.ID, err)
		return status.Errorf(codes.Internal, "failed to update message")
	}

	if res.RowsAffected() == 0 {
		return status.Errorf(codes.NotFound, "message not found")
	}

	return nil
}
