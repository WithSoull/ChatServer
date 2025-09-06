package chat_participant

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/ChatServer/internal/client/db"
	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/jackc/pgconn"
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return domainerrors.ErrUserAlreadyInChat
			case "23503": // foreign_key_violation
				return domainerrors.ErrChatNotFound
			case "23514": // check_violation
				return domainerrors.ErrInvalidRole
			}
		}
		return err
	}

	return nil
}
