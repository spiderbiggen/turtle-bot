package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/patrickmn/go-cache"
	cronLib "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"
	"turtle-bot/internal/command"
	kitsuApi "turtle-bot/internal/kitsu"
	nyaaApi "turtle-bot/internal/nyaa"
	"turtle-bot/internal/riot"
	"turtle-bot/internal/storage/couch"
	"turtle-bot/internal/storage/postgres"
	tenorApi "turtle-bot/internal/tenor"
	"turtle-bot/internal/worker"
)

var (
	commandHandlers    = make(map[string]command.Handler)
	componentHandlers  = make(map[string]command.Handler)
	registeredCommands []*discordgo.ApplicationCommand
)

var (
	logLevel string
)

func init() {
	flag.StringVar(&logLevel, "level", log.DebugLevel.String(), "Set log level, one of: trace, debug, info, warn, error, fatal, panic")
	rand.Seed(time.Now().Unix())
}

func main() {
	flag.Parse()
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetLevel(level)

	hostname, _ := os.Hostname()
	log.Infof("Starting Weeb Bot on %s", hostname)

	cron := cronLib.New()
	defer cron.Stop()

	log.Debugln("Migrating database...")
	db := postgres.New()
	couchdb := couch.New()
	kitsu := kitsuApi.New()
	nyaa := nyaaApi.New()
	tenor := tenorApi.New(os.Getenv("TENOR_KEY"))
	client := riot.New(os.Getenv("RIOT_KEY"))
	memCache := cache.New(5*time.Minute, 10*time.Minute)
	migrateDatabases(db, couchdb)

	d, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("TOKEN")))
	if err != nil {
		log.Fatal(err)
	}
	d.AddHandler(readyHandler(cron, db, couchdb, client, kitsu, nyaa, tenor, memCache))
	d.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	// Open a websocket connection to Discord and begin listening.
	err = d.Open()
	if err != nil {
		log.Fatal("Error opening discord connection,", err)
	}
	defer func() { _ = d.Close() }()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
}

func registerCommands(s *discordgo.Session, fs ...command.InteractionHandler) {
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
		if f == nil {
			continue
		}
		//cmd, h, c := f()
		cmd := f.Command()
		var v *discordgo.ApplicationCommand
		for _, ac := range commands {
			if ac.Name == cmd.Name && ac.Description == cmd.Description && len(ac.Options) == len(cmd.Options) {
				eq := true
				for i, option := range ac.Options {
					eq = eq && option == cmd.Options[i]
				}
				if eq {
					v = ac
					break
				}
			}
		}
		if v == nil {
			v, err = s.ApplicationCommandCreate(id, "", cmd)
			if err != nil {
				log.Errorf("Cannot create '%v' command: %v", cmd.Name, err)
				continue
			}
		}

		commandHandlers[f.InteractionID()] = f.HandleInteraction
		if v, ok := f.(command.ComponentHandler); ok {
			componentHandlers[v.ComponentID()] = v.HandleComponent
		}
		registeredCommands = append(registeredCommands, v)
		p = append(p, cmd.Name)
	}

	go func() {
		rem := filterRegisteredCommands(commands, p)
		for _, ac := range rem {
			if err := s.ApplicationCommandDelete(id, "", ac.ID); err != nil {
				log.Warnf("Failed to delete command: %v", err)
			}
		}
	}()

	log.Debugf("Started bot with registered commands: %s", strings.Join(p, ", "))
}

func filterRegisteredCommands(commands []*discordgo.ApplicationCommand, registeredCommands []string) []*discordgo.ApplicationCommand {
	var r []*discordgo.ApplicationCommand
o:
	for _, ac := range commands {
		for _, c := range registeredCommands {
			if ac.Name == c {
				continue o
			}
		}
		r = append(r, ac)
	}
	return r
}

func readyHandler(cron *cronLib.Cron, db *postgres.Client, couch *couch.Client, client *riot.Client, kitsu *kitsuApi.Client, nyaa *nyaaApi.Client, tenor *tenorApi.Client, memCache *cache.Cache) func(s *discordgo.Session, i *discordgo.Ready) {
	return func(s *discordgo.Session, i *discordgo.Ready) {
		matchQueue := worker.NewMatchWorker(client, couch)
		// Register commands if discord is ready
		registerCommands(s,
			&command.Apex{Client: tenor, Cache: memCache},
			&command.Play{Client: tenor, Cache: memCache},
			&command.Hurry{Client: tenor, Cache: memCache},
			&command.Morb{Client: tenor, Cache: memCache},
			&command.Sleep{Client: tenor, Cache: memCache},
			&command.RiotGroup{Api: client, Db: db, Couch: couch, Queue: matchQueue},
			command.AnimeGroup(kitsu, db),
		)

		var err error
		nyaa := worker.NyaaCheck(db, kitsu, nyaa)
		_, err = cron.AddFunc("*/30 * * * *", func() {
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
		}
		cron.Start()
	}
}

func migrateDatabases(db *postgres.Client, couchdb *couch.Client) {
	go func() {
		if err := db.Migrate(); err != nil {
			log.Errorf("Error migrating database: %v", err)
		} else {
			log.Debugf("Finished migration")
		}
	}()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := couchdb.Init(ctx); err != nil {
			log.Errorf("Error initializing database: %v", err)
		}
	}()
}
