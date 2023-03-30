package migration

import (
	"context"
	"github.com/jmoiron/sqlx"
)

func init() {
	Register(&SubscribeToTitle004{})
}

type SubscribeToTitle004 struct{}

func (t *SubscribeToTitle004) Index() int { return 4 }

func (t *SubscribeToTitle004) Name() string { return "change to substring subscriptions" }

func (t *SubscribeToTitle004) Up(ctx context.Context, tx *sqlx.Tx) error {
	sql := `
ALTER TABLE anime_has_subscriptions
    ADD COLUMN substring VARCHAR(255) DEFAULT NULL,
    DROP CONSTRAINT IF EXISTS anime_has_subscriptions_pkey;

UPDATE anime_has_subscriptions
SET substring = (SELECT query_title FROM anime WHERE anime.id = anime_has_subscriptions.anime_id);

ALTER TABLE anime_has_subscriptions
    ALTER COLUMN substring DROP DEFAULT,
    ALTER COLUMN substring SET NOT NULL,
    ALTER COLUMN anime_id DROP NOT NULL,
    ADD PRIMARY KEY (substring, channel_id, guild_id);
`
	_, err := tx.ExecContext(ctx, sql)
	return err
}

func (t *SubscribeToTitle004) Down(ctx context.Context, tx *sqlx.Tx) error {
	sql := `
ALTER TABLE anime_has_subscriptions
    DROP CONSTRAINT IF EXISTS anime_has_subscriptions_pkey,
    DROP COLUMN substring,
    ALTER COLUMN anime_id SET NOT NULL,
    ADD PRIMARY KEY (anime_id, channel_id, guild_id);`
	_, err := tx.ExecContext(ctx, sql)
	return err
}
