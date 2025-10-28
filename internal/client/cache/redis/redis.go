package redis

import (
	"context"
	"strconv"

	"github.com/WithSoull/ChatServer/internal/client/cache"
	"github.com/WithSoull/ChatServer/internal/config"
	"github.com/WithSoull/platform_common/pkg/logger"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

const (
	setKey = "users:active"
)

type handler func(ctx context.Context, conn redis.Conn) error

type usersIDsCacheClient struct {
	pool   *redis.Pool
	setKey string
}

func NewClient(pool *redis.Pool) cache.UsersIDsCacheClient {
	return &usersIDsCacheClient{
		pool:   pool,
		setKey: setKey,
	}
}

func (c *usersIDsCacheClient) Add(ctx context.Context, userID int64) error {
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		_, redisErr := conn.Do("SADD", c.setKey, strconv.FormatInt(userID, 10))
		return redisErr
	})
	if err != nil {
		logger.Debug(ctx, "Failed to add user to redis", zap.Int64("userID", userID))
		return err
	}
	return nil
}

func (c *usersIDsCacheClient) Remove(ctx context.Context, userID int64) error {
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		_, redisErr := conn.Do("SREM", c.setKey, strconv.FormatInt(userID, 10))
		return redisErr
	})
	if err != nil {
		logger.Error(ctx, "Failed to remove user to redis", zap.Int64("userID", userID))
		return err
	}
	return nil
}

func (c *usersIDsCacheClient) Exist(ctx context.Context, userID int64) (bool, error) {
	var exists bool
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		reply, redisErr := redis.Int(conn.Do("SISMEMBER", c.setKey, strconv.FormatInt(userID, 10)))
		if redisErr != nil {
			return redisErr
		}
		exists = reply == 1
		return nil
	})
	if err != nil {
		logger.Error(ctx, "Failed to check user exist in redis", zap.Int64("userID", userID))
		return false, err
	}
	return exists, nil
}

func (c *usersIDsCacheClient) execute(ctx context.Context, handler handler) error {
	conn, err := c.getConnect(ctx)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			logger.Error(ctx, "failed to close redis connection with error", zap.Error(err))
		}
	}()
	return handler(ctx, conn)
}

func (c *usersIDsCacheClient) getConnect(ctx context.Context) (redis.Conn, error) {
	getConnTimeoutCtx, cancel := context.WithTimeout(ctx, config.AppConfig().Redis.ConnTimeout())
	defer cancel()

	conn, err := c.pool.GetContext(getConnTimeoutCtx)
	if err != nil {
		logger.Error(ctx, "failed to get redis connection", zap.Error(err))

		_ = conn.Close()
		return nil, err
	}

	return conn, nil
}
