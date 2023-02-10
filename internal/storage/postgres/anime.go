package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Anime struct {
	ID             string       `db:"id"`
	CanonicalTitle string       `db:"canonical_title"`
	QueryTitle     string       `db:"query_title"`
	ImageURL       string       `db:"image_url"`
	CreatedAt      sql.NullTime `db:"created_at"`
}

type AnimeSubscription struct {
	AnimeID   string `db:"anime_id"`
	GuildID   string `db:"guild_id"`
	ChannelID string `db:"channel_id"`
}

type AnimeWithSubscriptions struct {
	Anime
	Subscriptions []AnimeSubscription
}

func (c *Client) InsertAnime(ctx context.Context, anime Anime) error {
	conn, err := c.Connection()
	if err != nil {
		return err
	}
	tx, err := conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareNamed(`
		INSERT INTO anime (id, canonical_title, query_title, image_url, created_at)
		VALUES (:id, :canonical_title, :query_title, :image_url, :created_at)
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		return err
	}
	if _, err = stmt.ExecContext(ctx, anime); err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}

func (c *Client) InsertAnimeSubscription(ctx context.Context, sub AnimeSubscription) error {
	conn, err := c.Connection()
	if err != nil {
		return err
	}
	tx, err := conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(`
		INSERT INTO anime_has_subscriptions (anime_id, guild_id, channel_id)
		VALUES (:anime_id, :guild_id, :channel_id)
		ON CONFLICT (anime_id, guild_id, channel_id) DO NOTHING
	`)

	if err != nil {
		return err
	}
	if _, err = stmt.ExecContext(ctx, sub); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}

func (c *Client) GetSubscriptions(ctx context.Context, queryTitle string) (withSubscriptions AnimeWithSubscriptions, err error) {
	var conn *sqlx.DB
	if conn, err = c.Connection(); err != nil {
		return
	}
	if err = conn.GetContext(ctx, &withSubscriptions.Anime, "SELECT * FROM anime WHERE query_title ILIKE $1 ORDER BY created_at DESC LIMIT 1", queryTitle); err != nil {
		err = fmt.Errorf("get anime: %w", err)
		return
	}

	if err = conn.SelectContext(ctx, &withSubscriptions.Subscriptions, "SELECT * FROM anime_has_subscriptions WHERE anime_id = $1", withSubscriptions.Anime.ID); err != nil {
		err = fmt.Errorf("select subs for %s: %w", withSubscriptions.Anime.ID, err)
	}
	return
}
