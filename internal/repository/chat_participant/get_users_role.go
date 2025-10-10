package chat_participant

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/platform_common/pkg/client/db"
	domainerrors "github.com/WithSoull/ChatServer/internal/errors/domain"
	"github.com/WithSoull/ChatServer/internal/model"
	"github.com/WithSoull/ChatServer/internal/repository/converter"
	rmodel "github.com/WithSoull/ChatServer/internal/repository/model"
	"github.com/jackc/pgx/v4"
)

// GetUserRole cant define who from user and chat does not exist
// that's why you should use checkUserRole from service layer to
// define this 2 corner cases
// To define the difference you need to select in chats table,
// so chat_participant table does not have access to another table
func (r *chatParticipantRepo) GetUserRole(ctx context.Context, chatID, userID int64) (model.Role, error) {
	builder := sq.Select(RoleColumn).
		From(TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{ChatIDColumn: chatID, UserIDColumn: userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return model.ROLE_USER, err
	}

	q := db.Query{
		Name:     "chat_participant_repository:GetUserRole",
		QueryRaw: query,
	}

	var role rmodel.Role
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ROLE_USER, domainerrors.ErrCantDefineWhoDoesNotExistUserOrChat
		}

		return model.ROLE_USER, err
	}

	return converter.FromRepoToModelRole(role), nil
}
