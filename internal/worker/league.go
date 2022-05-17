package worker

import (
	"context"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"weeb_bot/internal/riot"
	"weeb_bot/internal/storage"
)

var (
	matchSet          = make(map[string]interface{})
	retrievedMatchSet = make(map[string]interface{})
)

const channelId = "843990289374511114"

func MatchChecker(db *storage.Client, client *riot.Client) Worker {
	return func(ctx context.Context, discord *discordgo.Session) {
		summoners, err := db.GetSummoners(ctx)
		if err != nil {
			log.Errorf("Failed to get summoners from db: %v", err)
			return
		}
		downloadMatches(ctx, client, summoners)
	}
}

func downloadMatches(ctx context.Context, client *riot.Client, summoners []*riot.Summoner) {
	if len(summoners) == 0 {
		return
	}
	for _, summoner := range summoners {
		ids, err := client.MatchIds(ctx, riot.EUW1, summoner.Puuid, &riot.MatchIdsOptions{Count: 100})
		if err != nil {
			log.Warnf("Failed to get match info for %v", summoner)
		} else {
			for _, id := range ids {
				matchSet[id] = nil
			}
		}
	}
	r := make([]string, 0, 5)
	for id := range matchSet {
		if _, ok := retrievedMatchSet[id]; !ok {
			r = append(r, id)
		}
	}
	// TODO get batch matches
}

func updateSummoners(ctx context.Context, client *riot.Client, summoners []*riot.Summoner) []*riot.Summoner {
	r := make([]*riot.Summoner, 0, len(summoners))
	for _, summoner := range summoners {
		s, err := client.SummonerByPuuid(ctx, riot.EUW1, summoner.Puuid)
		if err != nil {
			log.Warnf("Failed to update summoner info for %v", summoner)
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
