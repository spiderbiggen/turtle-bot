package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
	kitsuApi "turtle-bot/internal/kitsu"
	nyaaApi "turtle-bot/internal/nyaa"
	"turtle-bot/internal/storage/postgres"
	"turtle-bot/internal/worker"
)

func main() {
	log.SetLevel(log.TraceLevel)
	db := postgres.New()
	kitsu := kitsuApi.New()
	nyaa := nyaaApi.New()

	d, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("TOKEN")))
	if err != nil {
		log.Fatal(err)
	}
	w := worker.NyaaCheck(db, kitsu, nyaa, time.Now().Add(-24*time.Hour))
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	w(ctx, d)
}
