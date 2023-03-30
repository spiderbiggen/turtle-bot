package postgres

import (
	"context"
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type Anime struct {
	ID             string       `db:"id"`
	CanonicalTitle string       `db:"canonical_title"`
	QueryTitle     string       `db:"query_title"`
	ImageURL       string       `db:"image_url"`
	CreatedAt      sql.NullTime `db:"created_at"`
}

type AnimeSubscription struct {
	AnimeID   sql.NullString `db:"anime_id"`
	Substring string         `db:"substring"`
	GuildID   string         `db:"guild_id"`
	ChannelID string         `db:"channel_id"`
}

type SubscriptionWithAnime struct {
	Subscription AnimeSubscription
	Anime        *Anime
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
		INSERT INTO anime_has_subscriptions (substring, guild_id, channel_id)
		VALUES (:substring, :guild_id, :channel_id)
		ON CONFLICT (substring, guild_id, channel_id) DO NOTHING
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

func (c *Client) GetSubscriptions(ctx context.Context, queryTitle string) ([]SubscriptionWithAnime, error) {
	conn, err := c.Connection()
	if err != nil {
		return nil, err
	}
	var subscriptions []AnimeSubscription
	err = conn.SelectContext(ctx, &subscriptions, "SELECT * FROM anime_has_subscriptions WHERE $1 ILIKE substring ", queryTitle)
	if err != nil {
		return nil, err
	}
	subscriptionsWithAnime := make([]SubscriptionWithAnime, len(subscriptions))
	for i, subscription := range subscriptions {
		subscriptionsWithAnime[i].Subscription = subscription
		if subscription.AnimeID.Valid {
			if err = conn.GetContext(ctx, &subscriptionsWithAnime[i].Anime, "SELECT * FROM anime WHERE id = $1", subscription.AnimeID.String); err != nil {
				log.Errorf("get anime: %v", err)
			}
		}
	}

	return subscriptionsWithAnime, nil
}
