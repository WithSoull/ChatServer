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
		log.Printf("failed to create message in chat %d: %v", rmsg.ChatID, err)
		return 0, status.Errorf(codes.Internal, "failed to create message")
	}

	return msgID, nil
}
