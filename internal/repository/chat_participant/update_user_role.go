package chat_participant

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/ChatServer/internal/client/db"
	"github.com/WithSoull/ChatServer/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		log.Printf("failed to update user %d role in chat %d: %v", userID, chatID, err)
		return status.Errorf(codes.Internal, "failed to update user role")
	}

	if res.RowsAffected() == 0 {
		return status.Errorf(codes.NotFound, "user not found in chat")
	}

	return nil
}
