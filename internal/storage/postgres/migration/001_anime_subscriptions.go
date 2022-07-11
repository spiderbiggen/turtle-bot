package migration

import (
	"context"
	"github.com/jmoiron/sqlx"
)

func init() {
	Register(&AnimeSubscriptions001{})
}

type AnimeSubscriptions001 struct{}

func (t *AnimeSubscriptions001) Index() int { return 1 }

func (t *AnimeSubscriptions001) Name() string { return "anime subscriptions" }

func (t *AnimeSubscriptions001) Up(ctx context.Context, tx *sqlx.Tx) error {
	sql := `
CREATE TABLE anime (
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP,
    canonical_title TEXT NOT NULL,
    query_title TEXT NOT NULL,
    image_url TEXT NOT NULL,
    UNIQUE (query_title)
);

CREATE TABLE anime_has_subscriptions (
  anime_id TEXT,
  guild_id TEXT,
  channel_id TEXT,
  PRIMARY KEY (anime_id, guild_id, channel_id)            
);`
	_, err := tx.ExecContext(ctx, sql)
	return err
}

func (t *AnimeSubscriptions001) Down(ctx context.Context, tx *sqlx.Tx) error {
	sql := `
DROP TABLE league_user;
DROP TABLE discord_user_has_league_user;`
	_, err := tx.ExecContext(ctx, sql)
	return err
}
