package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func (c *Client) PipeHSet(ctx context.Context, key string, value any) (int64, error) {
	return c.Pipe().HSet(ctx, key, value).Result()
}

func (c *Client) PipeSet(ctx context.Context, key string, value any, ttl time.Duration) (string, error) {
	return c.Pipe().Set(ctx, key, value, ttl).Result()
}

func (c *Client) PipeExec(ctx context.Context) ([]redis.Cmder, error) {
	return c.Pipe().Exec(ctx)
}
