package postgres

import (
	"context"
	"time"
	"weeb_bot/internal/riot"
)

type summoner struct {
	ID            string    `db:"id"`
	AccountId     string    `db:"account_id"`
	Puuid         string    `db:"puuid"`
	SummonerName  string    `db:"summoner_name"`
	RevisionDate  time.Time `db:"revision_date"`
	ProfileIconId int       `db:"profile_icon_id"`
	SummonerLevel int       `db:"summoner_level"`
}

func entityToSummoner(s summoner) *riot.Summoner {
	return &riot.Summoner{
		Id:            s.ID,
		AccountId:     s.AccountId,
		Puuid:         s.Puuid,
		SummonerName:  s.SummonerName,
		ProfileIconId: s.ProfileIconId,
		RevisionDate:  s.RevisionDate.UnixMilli(),
		SummonerLevel: s.SummonerLevel,
	}
}

func entityFromSummoner(r *riot.Summoner) *summoner {
	return &summoner{
		ID:            r.Id,
		AccountId:     r.AccountId,
		Puuid:         r.Puuid,
		SummonerName:  r.SummonerName,
		RevisionDate:  time.UnixMilli(r.RevisionDate),
		ProfileIconId: r.ProfileIconId,
		SummonerLevel: r.SummonerLevel,
	}
}

type discordUserHasLeagueUser struct {
	DiscordID string `db:"discord_id"`
	LeagueID  string `db:"league_id"`
}

func (c *Client) InsertDiscordSummoner(ctx context.Context, userID string, riotSummoner *riot.Summoner) error {
	conn, err := c.Connection()
	if err != nil {
		return err
	}
	tx, err := conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareNamed(`
		INSERT INTO league_user (id, account_id, puuid, profile_icon_id, summoner_name, summoner_level, revision_date)
		VALUES (:id, :account_id, :puuid, :profile_icon_id, :summoner_name, :summoner_level, :revision_date)
		ON CONFLICT (id) DO UPDATE SET account_id      = :account_id,
								  puuid           = :puuid,
								  profile_icon_id = :profile_icon_id,
								  summoner_name   = :summoner_name,
								  summoner_level  = :summoner_level,
								  revision_date   = :revision_date
	`)
	if err != nil {
		return err
	}
	if _, err = stmt.ExecContext(ctx, entityFromSummoner(riotSummoner)); err != nil {
		_ = tx.Rollback()
		return err
	}
	stmt, err = tx.PrepareNamed(`
		INSERT INTO discord_user_has_league_user (discord_id, league_id)
		VALUES (:discord_id, :league_id)
		ON CONFLICT (discord_id) DO UPDATE SET league_id = :league_id
	`)
	if err != nil {
		return err
	}
	if _, err = stmt.ExecContext(ctx, discordUserHasLeagueUser{DiscordID: userID, LeagueID: riotSummoner.Id}); err != nil {
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

func (c *Client) GetSummoners(ctx context.Context) ([]*riot.Summoner, error) {
	conn, err := c.Connection()
	if err != nil {
		return nil, err
	}
	var s []summoner
	if err = conn.SelectContext(ctx, &s, "SELECT * FROM league_user"); err != nil {
		return nil, err
	}
	r := make([]*riot.Summoner, len(s))
	for i, s2 := range s {
		r[i] = entityToSummoner(s2)
	}
	return r, nil
}
