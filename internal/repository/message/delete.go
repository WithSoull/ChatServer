package message

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/ChatServer/internal/client/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *messageRepo) Delete(ctx context.Context, messageID int64) error {
	builder := sq.Delete(TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{IDColumn: messageID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "message_repository:Delete",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("failed to delete message %d: %v", messageID, err)
		return status.Errorf(codes.Internal, "failed to delete message")
	}

	if res.RowsAffected() == 0 {
		return status.Errorf(codes.NotFound, "message not found")
	}

	return nil
}
