package worker

import (
	"context"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"weeb_bot/internal/riot"
)

var (
	summoners         []*riot.Summoner
	summonerNames     = [...]string{"TF Spiderbiggen", "TF Jellander", "TF Ralph", "TF Santa", "TF Super Sas"}
	matchSet          = make(map[string]interface{})
	retrievedMatchSet = make(map[string]interface{})
)

const channelId = "843990289374511114"

func MatchChecker(client *riot.Client) Worker {
	return func(ctx context.Context, discord *discordgo.Session) {
		if len(summoners) != len(summonerNames) {
			s := getSummoners(ctx, client)
			summoners = s
		}
	}
}

func getSummoners(ctx context.Context, c *riot.Client) []*riot.Summoner {
	r := make([]*riot.Summoner, 0, len(summoners))
	for _, name := range summonerNames {
		s, err := c.SummonerByName(ctx, riot.EUW1, name)
		if err != nil {
			log.Warnf("Failed to get summoner info for %s", name)
		} else {
			r = append(r, s)
		}
	}
	return r
}

func getMatches(ctx context.Context, c *riot.Client, summoners []*riot.Summoner) {
	if len(summoners) == 0 {
		return
	}

}
