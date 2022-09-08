package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
	"turtle-bot/internal/storage/postgres/migration"
)

type Client struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	db       *sqlx.DB
}

func New() *Client {
	return &Client{
		Username: os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASS"),
		Host:     os.Getenv("PG_HOST"),
		Port:     os.Getenv("PG_PORT"),
		Database: os.Getenv("PG_DATABASE"),
	}
}

func (c *Client) Connection() (*sqlx.DB, error) {
	var err error
	if c.db == nil {
		c.db, err = sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", c.Username, c.Database, c.Password, c.Host, c.Port))
		if err != nil {
			return nil, err
		}
	}
	return c.db, nil
}

func (c *Client) Migrate() error {
	db, err := c.Connection()
	if err != nil {
		return err
	}
	err = migration.Up(context.Background(), db)
	if err != nil {
		return err
	}
	return nil
}
