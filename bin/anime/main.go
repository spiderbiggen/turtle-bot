package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
	animeApi "turtle-bot/internal/anime"
	kitsuApi "turtle-bot/internal/kitsu"
	"turtle-bot/internal/storage/postgres"
	"turtle-bot/internal/worker"
)

func main() {
	log.SetLevel(log.TraceLevel)
	db := postgres.New()
	kitsu := kitsuApi.New()
	anime := animeApi.New()

	d, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("TOKEN")))
	if err != nil {
		log.Fatal(err)
	}
	w := worker.NewTorrent(db, kitsu, anime, time.Now().Add(-24*time.Hour))
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	w(ctx, d)
}
