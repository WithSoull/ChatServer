package message

import (
	"context"
	"database/sql"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/ChatServer/internal/client/db"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/ChatServer/internal/repository/converter"
	rmodel "github.com/WithSoull/ChatServer/internal/repository/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		if errors.Is(err, sql.ErrNoRows) {
			return model.Message{}, status.Errorf(codes.NotFound, "message not found")
		}
		log.Printf("failed to get message %d: %v", messageID, err)
		return model.Message{}, status.Errorf(codes.Internal, "failed to get message")
	}

	return converter.FromRepoToModelMessage(rmsg), nil
}
