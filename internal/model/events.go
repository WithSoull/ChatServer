package model

import "time"

type UserCreatedEvent struct {
	UserID    int64
	CreatedAt *time.Time
}

type UserDeletedEvent struct {
	UserID    int64
	DeletedAt *time.Time
}
