package chat_participant

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
)

func (r *chatParticipantRepo) GetUserRole(ctx context.Context, chatID, userID int64) (bool, model.Role) {
	builder := sq.Select(RoleColumn).
		From(TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{ChatIDColumn: chatID, UserIDColumn: userID})

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query for getUserRole (chatID=%d userID=%d): %v", chatID, userID, err)
		return false, model.ROLE_USER
	}

	q := db.Query{
		Name:     "chat_participant_repository:GetUserRole",
		QueryRaw: query,
	}

	var role rmodel.Role
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, model.ROLE_USER
		}
		log.Printf("failed to get user role (chatID=%d userID=%d): %v", chatID, userID, err)
		return false, model.ROLE_USER
	}

	return true, converter.FromRepoToModelRole(role)
}
