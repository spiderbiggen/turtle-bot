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
	"weeb_bot/internal/storage"
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
	db := storage.DefaultClient
	err := db.Migrate()
	if err != nil {
		log.Errorf("Error migrating database: %v", err)
	} else {
		log.Infof("Finished migration")
	}

	client := riot.New(os.Getenv("RIOT_KEY"))

	d, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("TOKEN")))
	if err != nil {
		log.Fatal(err)
	}
	d.AddHandler(readyHandler(cron, db, client))
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

	registerCommands(d,
		command.Sleep, command.Apex, command.Play, command.Hurry, command.Morbius, command.Morbin,
		command.RiotGroup(client, db),
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
}

func registerCommands(s *discordgo.Session, fs ...command.Factory) {
	commands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
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
			v, err = s.ApplicationCommandCreate(s.State.User.ID, "", c)
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

func readyHandler(cron *cronLib.Cron, db *storage.Client, client *riot.Client) func(s *discordgo.Session, i *discordgo.Ready) {
	return func(s *discordgo.Session, i *discordgo.Ready) {
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
			rito := worker.MatchChecker(db, client)
			_, err = cron.AddFunc("*/1 * * * *", func() {
				timeout, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancelFunc()
				rito(timeout, s)
			})
			if err != nil {
				log.Fatalln(err)
			}
		}
		cron.Start()
	}
}
