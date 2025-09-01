package chat

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

func (r *chatRepo) Get(ctx context.Context, chatID int64) (model.Chat, error) {
	builder := sq.Select(IDColumn, OwnerIDColumn, NameColumn, DescriptionColumn, CreatedAtColumn, UpdatedAtColumn).
		From(TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": chatID})

	query, args, err := builder.ToSql()
	if err != nil {
		return model.Chat{}, err
	}

	q := db.Query{
		Name:     "chat_repository:Get",
		QueryRaw: query,
	}

	var chat rmodel.Chat

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&chat.ID,
		&chat.OwnerID,
		&chat.Name,
		&chat.Description,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)
	if err != nil {
		log.Printf("failed to get chat (id=%d): %v", chatID, err)
		if errors.Is(err, sql.ErrNoRows) {
			return model.Chat{}, status.Errorf(codes.NotFound, "chat not found")
		}
		return model.Chat{}, status.Errorf(codes.Internal, "failed to get chat")
	}

	return converter.FromRepoToModelChat(chat, nil), nil
}
