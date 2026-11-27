package redis

import (
	"context"
	"errors"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

var ErrNotFound = errors.New("key not found")

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.RDB.Get(ctx, key).Result()
	switch {
	case errors.Is(err, goredis.Nil):
		return "", ErrNotFound
	case err != nil:
		return "", fmt.Errorf("redis get %q: %w", key, err)
	default:
		return val, nil
	}
}
