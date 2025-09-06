package domainerrors

import "errors"

var (
	ErrForbidden         = errors.New("action forbidden")
	ErrUserNotFound      = errors.New("user not found")
	ErrMessageNotFound   = errors.New("message not found")
	ErrChatNotFound      = errors.New("chat does not exist")
	ErrUserAlreadyInChat = errors.New("user is already in the chat")
	ErrInvalidRole       = errors.New("invalid role for participant")

	// This errors does not convert to grpc code
	// because its special intrenal errors for defining some corner cases
	ErrCantDefineWhoDoesNotExistUserOrChat = errors.New("cant define who does not exist: chat or user?")
)
