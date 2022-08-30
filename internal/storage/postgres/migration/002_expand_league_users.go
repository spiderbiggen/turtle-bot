package migration

import (
	"context"
	"github.com/jmoiron/sqlx"
)

func init() {
	Register(&ExpandLeagueUsers002{})
}

type ExpandLeagueUsers002 struct{}

func (t *ExpandLeagueUsers002) Index() int { return 2 }

func (t *ExpandLeagueUsers002) Name() string { return "expand league users" }

func (t *ExpandLeagueUsers002) Up(ctx context.Context, tx *sqlx.Tx) error {
	sql := `
ALTER TABLE league_user ADD COLUMN region INT NOT NULL DEFAULT 0;

ALTER TABLE discord_user_has_league_user 
    ALTER COLUMN league_id SET NOT NULL,
    ADD COLUMN channel_id TEXT;
`
	_, err := tx.ExecContext(ctx, sql)
	return err
}

func (t *ExpandLeagueUsers002) Down(ctx context.Context, tx *sqlx.Tx) error {
	sql := `
ALTER TABLE discord_user_has_league_user 
    DROP COLUMN IF EXISTS channel_id,
    ALTER COLUMN league_id DROP NOT NULL;

ALTER TABLE league_user DROP COLUMN IF EXISTS region;
`
	_, err := tx.ExecContext(ctx, sql)
	return err
}
