package chat_participant

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/ChatServer/internal/client/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *chatParticipantRepo) RemoveUser(ctx context.Context, chatID, userID int64) error {
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
		log.Printf("failed to remove user %d from chat %d: %v", userID, chatID, err)
		return status.Errorf(codes.Internal, "failed to remove user from chat")
	}

	if res.RowsAffected() == 0 {
		return status.Errorf(codes.NotFound, "user not found in chat")
	}

	return nil
}
