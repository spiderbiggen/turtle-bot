package worker

import (
	"context"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"sort"
	"weeb_bot/internal/riot"
	couch2 "weeb_bot/internal/storage/couch"
	"weeb_bot/internal/storage/models"
	postgres2 "weeb_bot/internal/storage/postgres"
)

type LeagueWorker struct {
	db    *postgres2.Client
	couch *couch2.Client
	*riot.Client
}

func MatchChecker(db *postgres2.Client, couch *couch2.Client, client *riot.Client) Worker {
	worker := LeagueWorker{
		db:     db,
		couch:  couch,
		Client: client,
	}
	return func(ctx context.Context, discord *discordgo.Session) {
		summoners, err := worker.db.GetSummoners(ctx)
		if err != nil {
			log.Errorf("Failed to get summoners from db: %v", err)
			return
		}
		worker.downloadMatches(ctx, summoners)
	}
}

func (w *LeagueWorker) downloadMatches(ctx context.Context, summoners []*models.RiotAccount) {
	if len(summoners) == 0 {
		return
	}

	ids, err := w.lastMatchesForSummoners(ctx, summoners)
	if err != nil {
		log.Errorf("Failed to get match ids for summoners: %v", err)
		return
	}
	iSize := len(ids)
	if iSize == 0 {
		return
	}

	l := 10
	if l > iSize {
		l = iSize
	}
	dst := make([]string, l)
	copy(dst, ids[:l])
	w.getMatches(ctx, dst)
}

func (w *LeagueWorker) lastMatchesForSummoners(ctx context.Context, summoners []*models.RiotAccount) ([]string, error) {
	c := len(summoners)
	if c == 0 {
		return nil, nil
	}

	matchSet := make(map[string]interface{})

	for _, summoner := range summoners {
		ids, err := w.MatchIds(ctx, riot.EUW1, summoner.Puuid, &riot.MatchIdsOptions{Count: 100})
		if err != nil {
			log.Warnf("Failed to get match info for %s", summoner.SummonerName)
			continue
		}

		for _, id := range ids {
			matchSet[id] = nil
		}
	}
	r := make([]string, 0, len(matchSet))
	for id := range matchSet {
		r = append(r, id)
	}
	// sort list from old to new
	sort.Strings(r)

	return w.couch.FilterMatchIds(ctx, r)
}

func (w *LeagueWorker) updateSummoners(ctx context.Context, summoners []*riot.Summoner) []*riot.Summoner {
	r := make([]*riot.Summoner, 0, len(summoners))
	for _, summoner := range summoners {
		s, err := w.SummonerByPuuid(ctx, riot.EUW1, summoner.Puuid)
		if err != nil {
			log.Warnf("Failed to update summoner info for %v", summoner)
		} else {
			r = append(r, s)
		}
	}
	return r
}

func (w *LeagueWorker) getMatches(ctx context.Context, matchIds []string) {
	if len(matchIds) == 0 {
		return
	}

	i := 0
	for _, id := range matchIds {
		match, err := w.Match(ctx, riot.EUW1, id)
		if err != nil {
			log.Errorf("Failed to retrieve match: %v", err)
			continue
		}
		if err := w.couch.AddMatch(ctx, match); err != nil {
			log.Warnf("Failed to add match to database: %v", err)
			continue
		}
		i++
	}
	if i == 0 {
		log.Errorf("Failed to save %d matches", len(matchIds))
		return
	}
	log.Infof("Succesfully saved %d matches", i)
}
