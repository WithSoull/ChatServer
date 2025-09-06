package chat_participant

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/ChatServer/internal/client/db"
)

func (r *chatParticipantRepo) GetUsers(ctx context.Context, chatID int64) ([]int64, error) {
	builder := sq.Select(UserIDColumn).
		From(TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{ChatIDColumn: chatID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "chat_participant_repository:GetUsers",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}
