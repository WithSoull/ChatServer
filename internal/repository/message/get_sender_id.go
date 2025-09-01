package message

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/ChatServer/internal/client/db"
)

func (r *messageRepo) GetMessageSenderID(ctx context.Context, messageID int64) (int64, error) {
	builder := sq.Select(SenderIDColumn).
		From(TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{IDColumn: messageID})

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query for GetMessageSenderID (messageID=%d): %v", messageID, err)
		return 0, err
	}

	q := db.Query{
		Name:     "message_repository:GetMessageSenderID",
		QueryRaw: query,
	}

	var senderID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&senderID)
	if err != nil {
		log.Printf("[Repository Layer] failed to execute query GetMessageSenderID (messageID=%d): %v", messageID, err)
		return 0, err
	}

	return senderID, nil
}
