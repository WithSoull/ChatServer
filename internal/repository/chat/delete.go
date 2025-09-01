package chat

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/ChatServer/internal/client/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		log.Printf("failed to delete chat (id=%d): %v", chatID, err)
		return status.Errorf(codes.Internal, "failed to delete chat")
	}

	if res.RowsAffected() == 0 {
		return status.Errorf(codes.NotFound, "chat not found")
	}

	return nil
}
