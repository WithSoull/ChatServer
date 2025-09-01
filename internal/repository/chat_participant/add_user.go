package chat_participant

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/ChatServer/internal/client/db"
	"github.com/WithSoull/ChatServer/internal/model"
)

func (r *chatParticipantRepo) AddUser(ctx context.Context, chatID, userID int64, role model.Role) error {
	now := time.Now()

	builder := sq.Insert(TableName).
		PlaceholderFormat(sq.Dollar).
		Columns(ChatIDColumn, UserIDColumn, RoleColumn, CreatedAtColumn, UpdatedAtColumn).
		Values(chatID, userID, role, now, now)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_participant_repository:AddUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		log.Printf("failed to add user %d to chat %d: %v", userID, chatID, err)
		return err
	}

	return nil
}
