package kafka

import "github.com/WithSoull/ChatServer/internal/model"

type UserCreatedDecoder interface {
	Decode(data []byte) (model.UserCreatedEvent, error)
}

type UserDeletedDecoder interface {
	Decode(data []byte) (model.UserDeletedEvent, error)
}
