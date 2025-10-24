package cache

import "context"

type UsersIDsCacheClient interface {
	Add(ctx context.Context, userID int64) error
	Remove(ctx context.Context, userID int64) error
	Exist(ctx context.Context, userID int64) (bool, error)
}
