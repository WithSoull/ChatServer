package chat_participant

import (
	"context"
	"strings"

	sq "github.com/Masterminds/squirrel"
	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	dmodel "github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/ChatServer/internal/repository/converter"
	rmodel "github.com/WithSoull/ChatServer/internal/repository/model"
	"github.com/WithSoull/platform_common/pkg/client/db"
)

func (r *chatParticipantRepo) RemoveUserFromChat(ctx context.Context, chatID, userID int64) error {
	builder := sq.Delete(TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{ChatIDColumn: chatID, UserIDColumn: userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "chat_participant_repository:RemoveUser",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return domainerrors.ErrUserNotFound(userID)
	}

	return nil
}

func (r *chatParticipantRepo) RemoveUserFromAllChats(ctx context.Context, userID int64) ([]dmodel.ChatParticipant, error) {
	builder := sq.Delete(TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{UserIDColumn: userID}).
		Suffix("RETURNING " +
			strings.Join([]string{
				ChatIDColumn,
				UserIDColumn,
				RoleColumn,
				CreatedAtColumn,
				UpdatedAtColumn,
			}, ", "),
		)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "chat_participant_repository:RemoveUserFromAllChats",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []dmodel.ChatParticipant
	for rows.Next() {
		var rr rmodel.ChatParticipant

		if err := rows.Scan(&rr.ChatID, &rr.UserID, &rr.Role, &rr.CreatedAt, &rr.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, converter.FromRepoToModelChatParticipant(rr))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
