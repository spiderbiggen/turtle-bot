package storage

type summoner struct {
	ID            string `db:"id"`
	AccountId     string `db:"account_id"`
	PlayerUuid    string `db:"player_uuid"`
	Name          string `db:"name"`
	ProfileIconId int    `db:"profile_icon_id"`
	RevisionDate  int64  `db:"revision_date"`
	SummonerLevel int    `db:"summoner_level"`
}
