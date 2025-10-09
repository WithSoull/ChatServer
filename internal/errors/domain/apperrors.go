package domainerrors

import (
	"github.com/WithSoull/platform_common/pkg/sys"
	"github.com/WithSoull/platform_common/pkg/sys/codes"
)

var (
	// Resource errors (NotFound)
	ErrUserNotFound    = sys.NewCommonError("user not found", codes.NotFound)
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
	ErrInvalidRole          = sys.NewCommonError("invalid role for participant", codes.InvalidArgument)
	ErrDuplicateParticipant = sys.NewCommonError("each participant can only be added once", codes.InvalidArgument)
	ErrNoChangesProvided    = sys.NewCommonError("no changes provided", codes.InvalidArgument)
	ErrEmptyMessageText     = sys.NewCommonError("message text cannot be empty", codes.InvalidArgument)

	// Conflict errors (AlreadyExists)
	ErrUserAlreadyInChat = sys.NewCommonError("user is already in the chat", codes.AlreadyExists)

	// Internal errors
	ErrFailedToCheckRole                   = sys.NewCommonError("failed to check user's role", codes.Internal)
	ErrCantDefineWhoDoesNotExistUserOrChat = sys.NewCommonError("cant define who does not exist: chat or user?", codes.Internal)
)
