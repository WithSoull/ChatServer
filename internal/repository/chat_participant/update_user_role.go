package chat_participant

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/platform_common/pkg/client/db"
	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
)

func (r *chatParticipantRepo) UpdateUserRole(ctx context.Context, chatID, userID int64, role model.Role) error {
	builder := sq.Update(TableName).
		PlaceholderFormat(sq.Dollar).
		Set(RoleColumn, role).
		Set(UpdatedAtColumn, time.Now()).
		Where(sq.Eq{ChatIDColumn: chatID, UserIDColumn: userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_participant_repository:UpdateUserRole",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return domainerrors.ErrUserNotFound
	}

	return nil
}
