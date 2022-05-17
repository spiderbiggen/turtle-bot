package storage

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
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

var schema = `
CREATE TABLE IF NOT EXISTS league_user (
    id TEXT PRIMARY KEY,
    account_id TEXT NOT NULL UNIQUE,
    puuid TEXT NOT NULL UNIQUE,
	summoner_name TEXT NOT NULL,
	summoner_level INTEGER NOT NULL,
	revision_date TIMESTAMP NOT NULL,
	profile_icon_id INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS discord_user_has_user (
  discord_id TEXT,
  league_id TEXT,
  PRIMARY KEY (discord_id, league_id),
  FOREIGN KEY (league_id) REFERENCES league_user (id) ON DELETE CASCADE
);
`

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
	_, err = db.Exec(schema)
	if err != nil {
		c.Enabled = false
		return err
	}
	return nil
}
