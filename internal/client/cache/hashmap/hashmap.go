package hashmap

import (
	"context"
	"sync"
)

// hashmapClient is a concurrency-safe set of user IDs.
type hashmapClient struct {
	mu      sync.RWMutex
	hashmap map[int64]struct{}
}

// NewClient returns a ready-to-use client with an initialized map.
func NewClient() *hashmapClient {
	return &hashmapClient{
		hashmap: make(map[int64]struct{}),
	}
}

// Add inserts the userID into the set.
func (c *hashmapClient) Add(ctx context.Context, userID int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	c.mu.Lock()
	if nil == c.hashmap {
		c.hashmap = make(map[int64]struct{})
	}
	c.hashmap[userID] = struct{}{}
	c.mu.Unlock()
	return nil
}

// Remove deletes the userID from the set (no-op if absent).
func (c *hashmapClient) Remove(ctx context.Context, userID int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	c.mu.Lock()
	delete(c.hashmap, userID)
	c.mu.Unlock()
	return nil
}

// Exist reports whether userID is in the set.
func (c *hashmapClient) Exist(ctx context.Context, userID int64) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	c.mu.RLock()
	_, ok := c.hashmap[userID]
	c.mu.RUnlock()
	return ok, nil
}
