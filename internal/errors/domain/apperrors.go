package domainerrors

import (
	"fmt"

	"github.com/WithSoull/platform_common/pkg/sys"
	"github.com/WithSoull/platform_common/pkg/sys/codes"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
)

var (
	// Resource errors (NotFound)
	ErrMessageNotFound = sys.NewCommonError("message not found", codes.NotFound)
	ErrChatNotFound    = sys.NewCommonError("chat does not exist", codes.NotFound)

	// Permission errors (PermissionDenied)
	ErrForbidden                = sys.NewCommonError("action forbidden", codes.PermissionDenied)
	ErrUserNotMember            = sys.NewCommonError("user are not a member of the chat", codes.PermissionDenied)
	ErrUserNoPermission         = sys.NewCommonError("user have no permission do it", codes.PermissionDenied)
	ErrOnlyAdminsAllowed        = sys.NewCommonError("only admins and owner has permission to make this action", codes.PermissionDenied)
	ErrOnlyOwnerAllowed         = sys.NewCommonError("only owner has permission to make this action", codes.PermissionDenied)
	ErrCannotRemoveHigherUser   = sys.NewCommonError("cannot remove a user with equal or higher role", codes.PermissionDenied)
	ErrCannotEditOthersMessages = sys.NewCommonError("you can edit only your own messages", codes.PermissionDenied)

	// Role management errors (PermissionDenied / InvalidArgument)
	ErrCannotAssignHigherRole     = sys.NewCommonError("cannot assign a role higher or equal than your own", codes.PermissionDenied)
	ErrCannotChangeHigherUserRole = sys.NewCommonError("cannot change role of a user with equal or higher role", codes.PermissionDenied)
	ErrCannotAssignOwnerRole      = sys.NewCommonError("cannot assign owner role to another member", codes.InvalidArgument)

	// Validation errors (InvalidArgument)
	ErrInvalidRole          = validate.NewValidationErrors("invalid role for participant")
	ErrDuplicateParticipant = validate.NewValidationErrors("each participant can only be added once")
	ErrNoChangesProvided    = validate.NewValidationErrors("no changes provided")
	ErrEmptyChatName        = validate.NewValidationErrors("name of the chat cannot be empty")
	ErrEmptyMessageText     = validate.NewValidationErrors("message text cannot be empty")

	// Conflict errors (AlreadyExists)
	ErrUserAlreadyInChat = sys.NewCommonError("user is already in the chat", codes.AlreadyExists)

	// Internal errors
	ErrFailedToCheckRole                   = sys.NewCommonError("failed to check user's role", codes.Internal)
	ErrCantDefineWhoDoesNotExistUserOrChat = sys.NewCommonError("cant define who does not exist: chat or user?", codes.Internal)
)

func ErrUserNotFound(userID int64) error {
	return sys.NewCommonError(fmt.Sprintf("user(ID=%d) not found", userID), codes.NotFound)
}
