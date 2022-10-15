package migration

import (
	"context"
	"github.com/jmoiron/sqlx"
)

func init() {
	Register(&ExpandLeagueUsers003{})
}

type ExpandLeagueUsers003 struct{}

func (t *ExpandLeagueUsers003) Index() int { return 3 }

func (t *ExpandLeagueUsers003) Name() string { return "remove anime unique query title" }

func (t *ExpandLeagueUsers003) Up(ctx context.Context, tx *sqlx.Tx) error {
	sql := `
ALTER TABLE anime DROP CONSTRAINT anime_query_title_key 
`
	_, err := tx.ExecContext(ctx, sql)
	return err
}

func (t *ExpandLeagueUsers003) Down(ctx context.Context, tx *sqlx.Tx) error {
	sql := `
ALTER TABLE anime ADD CONSTRAINT anime_query_title_key UNIQUE (query_title)
`
	_, err := tx.ExecContext(ctx, sql)
	return err
}
