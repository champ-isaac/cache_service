package redis

import "context"

func (c *Client) HSet(ctx context.Context, key string, value any) (int64, error) {
	return c.raw.HSet(ctx, key, value).Result()
}

func (c *Client) HGetAll(ctx context.Context, key string) (any, error) {
	return c.raw.HGetAll(ctx, key).Result()
}

func (c *Client) HDel(ctx context.Context, key string) (int64, error) {
	return c.raw.HDel(ctx, key).Result()
}

func (c *Client) HGetAllNScan(ctx context.Context, key string, dest any) error {
	return c.raw.HGetAll(ctx, key).Scan(dest)
}

func (c *Client) SMember(ctx context.Context, key string) (any, error) {
	return c.raw.SMembers(ctx, key).Result()
}
