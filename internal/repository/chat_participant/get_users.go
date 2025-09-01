package chat_participant

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/ChatServer/internal/client/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		log.Printf("failed to get users from chat %d: %v", chatID, err)
		return nil, status.Errorf(codes.Internal, "failed to get users from chat")
	}
	defer rows.Close()

	var userIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Printf("failed to scan user id from chat %d: %v", chatID, err)
			return nil, status.Errorf(codes.Internal, "failed to get users from chat")
		}
		userIDs = append(userIDs, id)
	}

	if err := rows.Err(); err != nil {
		log.Printf("iteration error while scanning users from chat %d: %v", chatID, err)
		return nil, status.Errorf(codes.Internal, "failed to get users from chat")
	}

	return userIDs, nil
}
