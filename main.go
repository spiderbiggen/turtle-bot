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
	animeApi "turtle-bot/internal/anime"
	"turtle-bot/internal/command"
	kitsuApi "turtle-bot/internal/kitsu"
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

type AppContext struct {
	DB       *postgres.Client
	Kitsu    *kitsuApi.Client
	Tenor    *tenorApi.Client
	Anime    *animeApi.Client
	MemCache *cache.Cache
	Cron     *cronLib.Cron
}

func main() {
	flag.Parse()
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetLevel(level)

	appContext := AppContext{
		DB:       postgres.New(),
		Kitsu:    kitsuApi.New(),
		Anime:    animeApi.New(),
		Tenor:    tenorApi.New(os.Getenv("TENOR_KEY")),
		Cron:     cronLib.New(),
		MemCache: cache.New(5*time.Minute, 10*time.Minute),
	}
	defer appContext.Cron.Stop()

	log.Debugln("Migrating database...")
	go migrateDatabases(appContext.DB)

	d, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("TOKEN")))
	if err != nil {
		log.Fatal(err)
	}
	d.AddHandler(readyHandler(appContext))
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
	hostname, _ := os.Hostname()
	log.Infof("Started Turtle Bot on %s", hostname)

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

func readyHandler(appContext AppContext) func(s *discordgo.Session, i *discordgo.Ready) {
	return func(s *discordgo.Session, i *discordgo.Ready) {
		// Register commands if discord is ready
		registerCommands(s,
			&command.Apex{Client: appContext.Tenor, Cache: appContext.MemCache},
			&command.Warzone{Client: appContext.Tenor, Cache: appContext.MemCache},
			&command.Play{Client: appContext.Tenor, Cache: appContext.MemCache},
			&command.Hurry{Client: appContext.Tenor, Cache: appContext.MemCache},
			&command.Morb{Client: appContext.Tenor, Cache: appContext.MemCache},
			&command.Sleep{Client: appContext.Tenor, Cache: appContext.MemCache},
			command.AnimeGroup(appContext.Kitsu, appContext.DB),
		)

		appContext.Cron.Start()
		nyaaWorker := worker.NewTorrent(appContext.DB, appContext.Kitsu, appContext.Anime)
		if err := nyaaWorker.Schedule(appContext.Cron, s); err != nil {
			log.Fatalln(err)
		}
		if err := nyaaWorker.Run(context.Background(), s); err != nil {
			log.Error(err)
		}
	}
}

func migrateDatabases(db *postgres.Client) {
	if err := db.Migrate(); err != nil {
		log.Errorf("Error migrating database: %v", err)
		panic(err)
	} else {
		log.Debugf("Finished migration")
	}
}
