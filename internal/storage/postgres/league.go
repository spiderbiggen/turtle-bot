package postgres

import (
	"context"
	"time"
	"weeb_bot/internal/riot"
	"weeb_bot/internal/storage/models"
)

type summoner struct {
	ID            string    `db:"id"`
	AccountId     string    `db:"account_id"`
	Puuid         string    `db:"puuid"`
	SummonerName  string    `db:"summoner_name"`
	RevisionDate  time.Time `db:"revision_date"`
	ProfileIconId int       `db:"profile_icon_id"`
	SummonerLevel int       `db:"summoner_level"`
	Region        uint8     `db:"region"`
}

func entityToRiotAccount(s summoner) models.RiotAccount {
	return models.RiotAccount{
		Summoner: riot.Summoner{
			Id:            s.ID,
			AccountId:     s.AccountId,
			Puuid:         s.Puuid,
			SummonerName:  s.SummonerName,
			ProfileIconId: s.ProfileIconId,
			RevisionDate:  s.RevisionDate.UnixMilli(),
			SummonerLevel: s.SummonerLevel,
		},
		Region: riot.Region(s.Region),
	}
}

func entityFromRiotAccount(acc models.RiotAccount) summoner {
	return summoner{
		ID:            acc.Id,
		AccountId:     acc.AccountId,
		Puuid:         acc.Puuid,
		SummonerName:  acc.SummonerName,
		RevisionDate:  time.UnixMilli(acc.RevisionDate),
		ProfileIconId: acc.ProfileIconId,
		SummonerLevel: acc.SummonerLevel,
		Region:        uint8(acc.Region),
	}
}

type discordUserHasLeagueUser struct {
	DiscordID string `db:"discord_id"`
	LeagueID  string `db:"league_id"`
	ChannelID string `db:"channel_id"`
}

func (c *Client) InsertDiscordSummoner(ctx context.Context, userID, channelID string, acc models.RiotAccount) error {
	conn, err := c.Connection()
	if err != nil {
		return err
	}
	tx, err := conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareNamed(`
		INSERT INTO league_user (id, account_id, puuid, profile_icon_id, summoner_name, summoner_level, revision_date, region)
		VALUES (:id, :account_id, :puuid, :profile_icon_id, :summoner_name, :summoner_level, :revision_date, :region)
		ON CONFLICT (id) DO UPDATE SET account_id      = :account_id,
								  puuid           = :puuid,
								  profile_icon_id = :profile_icon_id,
								  summoner_name   = :summoner_name,
								  summoner_level  = :summoner_level,
								  revision_date   = :revision_date,
								  region = :region
	`)
	if err != nil {
		return err
	}
	if _, err = stmt.ExecContext(ctx, entityFromRiotAccount(acc)); err != nil {
		_ = tx.Rollback()
		return err
	}
	stmt, err = tx.PrepareNamed(`
		INSERT INTO discord_user_has_league_user (discord_id, league_id, channel_id)
		VALUES (:discord_id, :league_id, :channel_id)
		ON CONFLICT (discord_id) DO UPDATE SET league_id = :league_id, channel_id = :channel_id
	`)
	if err != nil {
		return err
	}
	rel := discordUserHasLeagueUser{DiscordID: userID, LeagueID: acc.Id, ChannelID: channelID}
	if _, err = stmt.ExecContext(ctx, rel); err != nil {
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

func (c *Client) GetSummoners(ctx context.Context) ([]*models.RiotAccount, error) {
	conn, err := c.Connection()
	if err != nil {
		return nil, err
	}
	var s []summoner
	if err = conn.SelectContext(ctx, &s, "SELECT * FROM league_user"); err != nil {
		return nil, err
	}
	r := make([]*models.RiotAccount, len(s))
	for i, s2 := range s {
		acc := entityToRiotAccount(s2)
		r[i] = &acc
	}
	return r, nil
}

func (c *Client) GetDiscordSummoner(ctx context.Context, userID string) (channelID string, account *models.RiotAccount, err error) {
	conn, err := c.Connection()
	if err != nil {
		return
	}
	var rel struct {
		summoner
		discordUserHasLeagueUser
	}
	query := "SELECT * FROM discord_user_has_league_user du JOIN league_user lu ON lu.id = du.league_id WHERE discord_id = $1"
	if err = conn.GetContext(ctx, &rel, query, userID); err != nil {
		return
	}
	acc := entityToRiotAccount(rel.summoner)
	return rel.ChannelID, &acc, nil
}
