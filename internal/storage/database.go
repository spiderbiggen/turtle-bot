package storage

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

type Config struct {
	Username string
	Password string
	Host     string
	Database string
}

var schema = `
CREATE TABLE IF NOT EXISTS league_match (
    
);
`

var defaultConfig = &Config{
	Username: os.Getenv("PG_USER"),
	Password: os.Getenv("PG_PASS"),
	Host:     os.Getenv("PG_HOST"),
	Database: os.Getenv("PG_DATABASE"),
}

func (c *Config) connection(ctx context.Context) {
	_, _ = sqlx.ConnectContext(ctx, "postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s", c.Username, c.Database, c.Password, c.Host))
}
