package redis

import (
	"github.com/redis/go-redis/v9"
	"os"
)

type Client struct {
	Address  string
	Password string
	Database int
	db       *redis.Client
}

func New() *Client {
	return &Client{
		Address:  os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Database: 0,
	}
}

func (c *Client) Connection() *redis.Client {
	if c.db == nil {
		c.db = redis.NewClient(&redis.Options{
			Addr:     c.Address,
			Password: c.Password,
			DB:       c.Database,
		})
	}
	return c.db
}
