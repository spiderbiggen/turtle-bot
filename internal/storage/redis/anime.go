package redis

import (
	"context"
	"time"
)

const animeKey = "anime_check"

func (c *Client) SetLastAnimeSync(ctx context.Context, time time.Time) error {
	conn := c.Connection()
	return conn.Set(ctx, animeKey, time, 0).Err()
}

func (c *Client) GetLastAnimeSync(ctx context.Context) (time.Time, error) {
	conn := c.Connection()
	return conn.Get(ctx, animeKey).Time()
}
