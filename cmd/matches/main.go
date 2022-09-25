package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"sync"
	"time"
	"turtle-bot/internal/riot"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetLevel(log.TraceLevel)
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, "batch", false)
	defer cancel()
	api := riot.New(os.Getenv("RIOT_KEY"))

	summoner, err := api.SummonerByName(ctx, riot.EUW1, "TF Spiderbiggen")
	if err != nil {
		log.Fatal(err)
	}
	matchIds := make([]string, 0, 500)
	for pos := int32(0); ; {
		ids, _ := api.MatchIds(ctx, riot.EUW1, summoner.Puuid, &riot.MatchIdsOptions{Count: 100, Start: pos})
		matchIds = append(matchIds, ids...)
		pos += 100
		if len(ids) < 100 {
			break
		}
	}
	if len(matchIds) == 0 {
		log.Fatal("No matches")
	}
	wg := sync.WaitGroup{}
	for _, id := range matchIds {
		go func(id string) {
			wg.Add(1)
			defer wg.Done()
			_, err := api.Match(ctx, riot.EUW1, id)
			if err != nil {
				log.Error(err)
			}
		}(id)
	}
	wg.Wait()
}
