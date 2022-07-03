package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	cronLib "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"
	"weeb_bot/internal/command"
	"weeb_bot/internal/riot"
	"weeb_bot/internal/storage/couch"
	"weeb_bot/internal/storage/postgres"
	"weeb_bot/internal/worker"
)

var (
	commandHandlers    = make(map[string]command.Handler)
	registeredCommands []*discordgo.ApplicationCommand
)

func main() {
	rand.Seed(time.Now().Unix())
	log.SetLevel(log.TraceLevel)

	hostname, _ := os.Hostname()
	log.Infof("Starting Weeb Bot on %s", hostname)

	cron := cronLib.New()
	defer cron.Stop()

	log.Infoln("Migrating database...")
	db := postgres.New()
	err := db.Migrate()
	if err != nil {
		log.Errorf("Error migrating database: %v", err)
	} else {
		log.Infof("Finished migration")
	}
	couchdb := couch.New()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		err = couchdb.Init(ctx)
		if err != nil {
			log.Errorf("Error initializing database: %v", err)
		}
	}()

	client := riot.New(os.Getenv("RIOT_KEY"))

	d, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("TOKEN")))
	if err != nil {
		log.Fatal(err)
	}
	d.AddHandler(readyHandler(cron, db, couchdb, client))
	d.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	// Open a websocket connection to Discord and begin listening.
	err = d.Open()
	if err != nil {
		log.Fatal("Error opening discord connection,", err)
	}
	defer d.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
}

func registerCommands(s *discordgo.Session, fs ...command.Factory) {
	if s.State == nil || s.State.User == nil {
		log.Error("Missing user ID in session. Cannot register commands")
		return
	}
	id := s.State.User.ID
	commands, err := s.ApplicationCommands(id, "")
	if err != nil {
		log.Errorf("Failed to fetch current commands. %v", err)
		return
	}
	var p []string
	for _, f := range fs {
		c, h := f()
		var v *discordgo.ApplicationCommand
		for _, ac := range commands {
			if ac.Name == c.Name && ac.Description == c.Description && len(ac.Options) == len(c.Options) {
				same := true
				for i, option := range ac.Options {
					same = same && option == c.Options[i]
				}
				if same {
					v = ac
					break
				}
			}
		}
		if v == nil {
			v, err = s.ApplicationCommandCreate(id, "", c)
			if err != nil {
				log.Errorf("Cannot create '%v' command: %v", c.Name, err)
				continue
			}
		}

		commandHandlers[c.Name] = h
		registeredCommands = append(registeredCommands, v)
		p = append(p, c.Name)
	}
	log.Infof("Started bot with registered commands: %s.", strings.Join(p, ", "))
}

func readyHandler(cron *cronLib.Cron, db *postgres.Client, couch *couch.Client, client *riot.Client) func(s *discordgo.Session, i *discordgo.Ready) {
	return func(s *discordgo.Session, i *discordgo.Ready) {
		// Register commands if discord is ready
		registerCommands(s,
			command.Sleep, command.Apex, command.Play, command.Hurry,
			command.Morbius, command.Morbin, command.Morb,
			command.RiotGroup(client, db),
		)

		var err error
		nyaa := worker.NyaaCheck()
		_, err = cron.AddFunc("*/10 * * * *", func() {
			timeout, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancelFunc()
			nyaa(timeout, s)
		})
		if err != nil {
			log.Fatalln(err)
		}
		if db.Enabled {
			rito := worker.MatchChecker(db, couch, client)
			cmd := func() {
				timeout, cancelFunc := context.WithTimeout(context.Background(), 1*time.Minute)
				defer cancelFunc()
				rito(timeout, s)
			}
			_, err = cron.AddFunc("*/5 * * * *", cmd)
			if err != nil {
				log.Fatalln(err)
			}
			cmd()
		}
		cron.Start()
	}
}
