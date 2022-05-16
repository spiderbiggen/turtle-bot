package riot

import (
	"context"
	"encoding/json"
	"fmt"
)

type Summoner struct {
	Id            string `json:"id"`
	AccountId     string `json:"accountId"`
	Puuid         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

func (c *Client) SummonerByName(ctx context.Context, region Region, summonerName string) (*Summoner, error) {
	r, v := region.Realm()
	if !v {
		return nil, ErrRegionUnknown
	}
	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s", r, summonerName)
	resp, err := c.request(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s Summoner
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}
