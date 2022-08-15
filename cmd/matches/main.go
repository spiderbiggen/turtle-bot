package main

import (
	"context"
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
	"weeb_bot/internal/command"
	"weeb_bot/internal/riot"
	"weeb_bot/internal/storage/couch"
	"weeb_bot/internal/storage/postgres"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetLevel(log.DebugLevel)
	var err error
	couchdb := couch.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	conn, err := couchdb.Connection()
	if err != nil {
		log.Fatal(err)
	}

	name, err := randomName()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = conn.DestroyDB(context.Background(), name) }()
	if err := conn.CreateDB(ctx, name); err != nil {
		log.Fatal(err)
	}
	couchdb.Database = name
	if err := couchdb.Init(ctx); err != nil {
		log.Fatal(err)
	}
	g := command.RiotGroup{Api: riot.New(os.Getenv("RIOT_KEY")), Db: postgres.New(), Couch: couchdb}

	summoner, err := g.Api.SummonerByName(ctx, riot.EUW1, "TF Spiderbiggen")
	if err != nil {
		log.Fatal(err)
	}
	g.GetMatchHistory(summoner, riot.EUW1)
}

func randomName() (string, error) {
	var result string
	arr := make([]byte, 24)
	for len(result) == 0 || unicode.IsDigit(rune(result[0])) {
		if _, err := rand.Read(arr); err != nil {
			continue
		}
		result = strings.ToLower(base64.URLEncoding.EncodeToString(arr))
	}
	return result, nil
}
