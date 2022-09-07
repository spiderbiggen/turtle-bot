package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
	"turtle-bot/internal/riot"
	"turtle-bot/internal/stats"
	"turtle-bot/internal/storage/couch"
)

func main() {
	var err error
	couchdb := couch.New()
	api := riot.New(os.Getenv("RIOT_KEY"))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = couchdb.Init(ctx)
	if err != nil {
		log.Fatalf("Init failed: %v", err)
	}
	puuid := "jealcFNOjVdIZAiDIPcUlEGXeEa4tyK4voa59xJkRZF7GtpgTMC2stABJjFxRHB4dwPyiUcJEqfLYA"

	ids, err := api.MatchIds(ctx, riot.EUW1, puuid, &riot.MatchIdsOptions{Count: 1})
	if err != nil {
		log.Fatalf("MatchIds failed: %v", err)
	}
	if len(ids) == 0 {
		log.Fatal("No matches")
	}
	m, err := api.Match(ctx, riot.EUW1, ids[0])
	if err != nil {
		log.Fatalf("Match failed: %v", err)
	}
	participant := participantByPuuid(m.Info.Participants, puuid)
	if participant == nil {
		log.Fatalf("No participant for puuid %q", puuid)
	}

	avg, err := couchdb.GetQueueAverages(ctx, m.Info.QueueID, puuid)
	if err != nil {
		log.Fatalf("GetAverages failed: %v", err)
	}

	cmp, err := stats.FindComparable(participant, avg)
	if err != nil {
		log.Fatalf("FindComparable failed: %v", err)
	}
	log.Info("=====MAX=====")
	cmp2 := cmp.FilterMax()
	for ct, result := range cmp2 {
		log.Infof("%s, %+v, %.2f", ct, result, result.Average())
	}
	log.Info("=====MIN=====")
	cmp2 = cmp.FilterMin()
	for ct, result := range cmp2 {
		log.Infof("%s, %+v, %.2f", ct, result, result.Average())
	}

	log.Info("=====AV+=====")
	cmp2 = cmp.FilterAboveAverage()
	for ct, result := range cmp2 {
		log.Infof("%s, %+v, %.2f", ct, result, result.Average())
	}

	log.Info("=====AV-=====")
	cmp2 = cmp.FilterBelowAverage()
	for ct, result := range cmp2 {
		log.Infof("%s, %+v, %.2f", ct, result, result.Average())
	}
}

func participantByPuuid(participants []*riot.Participant, puuid string) *riot.Participant {
	for _, p := range participants {
		if p.Puuid == puuid {
			return p
		}
	}
	return nil
}
