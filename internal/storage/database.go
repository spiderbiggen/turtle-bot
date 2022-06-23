package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
	"weeb_bot/internal/storage/migration"
)

type Client struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	db       *sqlx.DB
	Enabled  bool
}

var (
	ErrNoConnection = errors.New("cannot connect to database")
)

var DefaultClient = &Client{
	Username: os.Getenv("PG_USER"),
	Password: os.Getenv("PG_PASS"),
	Host:     os.Getenv("PG_HOST"),
	Port:     os.Getenv("PG_PORT"),
	Database: os.Getenv("PG_DATABASE"),
	Enabled:  true,
}

func (c *Client) Connection() (*sqlx.DB, error) {
	if !c.Enabled {
		return nil, ErrNoConnection
	}
	var err error
	if c.db == nil {
		c.db, err = sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", c.Username, c.Database, c.Password, c.Host, c.Port))
		if err != nil {
			return nil, err
		}
		return c.db, nil
	}
	if err := c.db.Ping(); err != nil {
		return nil, err
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
		c.Enabled = false
		return err
	}
	return nil
}
