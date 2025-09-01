package chat

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

func (r *chatRepo) Create(ctx context.Context, chat model.Chat) (int64, error) {
	rChat := converter.FromModelToRepoChat(chat)
	now := time.Now()
	builder := sq.Insert(TableName).
		PlaceholderFormat(sq.Dollar).
		Columns(OwnerIDColumn, NameColumn, DescriptionColumn, CreatedAtColumn, UpdatedAtColumn).
		Values(rChat.OwnerID, rChat.Name, rChat.Description, now, now).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "chat_repository:Create",
		QueryRaw: query,
	}

	var chatID int64

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)

	if err != nil {
		log.Printf("[repo] failed to create chat: %v", err)
		return 0, status.Errorf(codes.Internal, "failed to create chat")
	}

	return chatID, nil
}
