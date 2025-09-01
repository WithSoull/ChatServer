package message

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/ChatServer/internal/client/db"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/ChatServer/internal/repository/converter"
	rmodel "github.com/WithSoull/ChatServer/internal/repository/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		log.Printf("failed to get messages from chat %d: %v", chatID, err)
		return nil, status.Errorf(codes.Internal, "failed to get messages")
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
			log.Printf("failed to scan message from chat %d: %v", chatID, err)
			return nil, status.Errorf(codes.Internal, "failed to get messages")
		}
		rmessages = append(rmessages, rmsg)
	}

	if err := rows.Err(); err != nil {
		log.Printf("iteration error while scanning messages from chat %d: %v", chatID, err)
		return nil, status.Errorf(codes.Internal, "failed to get messages")
	}
	log.Printf("[Repository Layer - message] rmessages=%+v", rmessages)
	log.Printf("[Repository Layer - message] messages=%+v", converter.FromRepoToModelMessages(rmessages))
	return converter.FromRepoToModelMessages(rmessages), nil
}
