package redis

import (
	"context"
	"time"
)

func (c *Client) Get(ctx context.Context, key string) (any, error) {
	return c.raw.Get(ctx, key).Result()
}

func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.raw.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.raw.Del(ctx, key).Err()
}

func (c *Client) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.raw.Scan(ctx, cursor, match, count).Result()
}

func (c *Client) FlushDb(ctx context.Context) error {
	return c.raw.FlushDB(ctx).Err()
}
