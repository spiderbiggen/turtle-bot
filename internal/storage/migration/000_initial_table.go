package migration

import (
	"context"
	"github.com/jmoiron/sqlx"
)

func init() {
	Register(&InitialTable001{})
}

type InitialTable001 struct{}

func (t *InitialTable001) Index() int { return 0 }

func (t *InitialTable001) Name() string { return "initial migration" }

func (t *InitialTable001) Up(ctx context.Context, tx *sqlx.Tx) error {
	sql := `
CREATE TABLE league_user (
    id TEXT PRIMARY KEY,
    account_id TEXT NOT NULL UNIQUE,
    puuid TEXT NOT NULL UNIQUE,
	summoner_name TEXT NOT NULL,
	summoner_level INTEGER NOT NULL,
	revision_date TIMESTAMP NOT NULL,
	profile_icon_id INTEGER NOT NULL
);

CREATE TABLE discord_user_has_league_user (
  discord_id TEXT PRIMARY KEY,
  league_id TEXT,
  FOREIGN KEY (league_id) REFERENCES league_user (id) ON DELETE CASCADE
);`
	_, err := tx.ExecContext(ctx, sql)
	return err
}

func (t *InitialTable001) Down(ctx context.Context, tx *sqlx.Tx) error {
	sql := `
DROP TABLE league_user;
DROP TABLE discord_user_has_league_user;`
	_, err := tx.ExecContext(ctx, sql)
	return err
}
