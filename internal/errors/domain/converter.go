package domainerrors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Return true if log is needed, and error ofc
func ToGRPCStatus(err error) (bool, error) {
	switch {
	case errors.Is(err, ErrForbidden):
		return false, status.Error(codes.PermissionDenied, "You are not allowed to perform this action")
	case errors.Is(err, ErrUserNotFound):
		return false, status.Error(codes.NotFound, "The specified user does not exist in this chat")
	case errors.Is(err, ErrMessageNotFound):
		return false, status.Error(codes.NotFound, "The message was not found")
	case errors.Is(err, ErrChatNotFound):
		return false, status.Error(codes.NotFound, "The requested chat does not exist")
	case errors.Is(err, ErrUserAlreadyInChat):
		return false, status.Error(codes.AlreadyExists, "The user already in chat")
	default:
		return true, status.Error(codes.Internal, "unknown internal error")
	}
}
