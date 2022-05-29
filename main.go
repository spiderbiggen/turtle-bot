package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"
	"weeb_bot/internal/command"
	"weeb_bot/internal/worker"
)

var (
	commandHandlers    = make(map[string]command.Handler)
	registeredCommands []*discordgo.ApplicationCommand
)

func main() {
	log.SetReportCaller(true)
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.TraceLevel)

	rand.Seed(time.Now().UnixNano())

	c := cron.New()
	defer c.Stop()
	hostname, _ := os.Hostname()
	log.Warnf("Starting Weeb Bot on %s", hostname)

	d, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("TOKEN")))
	if err != nil {
		log.Fatal(err)
	}
	d.AddHandler(func(s *discordgo.Session, i *discordgo.Ready) {
		check := worker.NyaaCheck()
		_, err := c.AddFunc("*/10 * * * *", func() {
			timeout, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancelFunc()
			check(timeout, s)
		})
		if err != nil {
			log.Fatalln(err)
		}
		c.Start()
	})
	d.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	// Open a websocket connection to Discord and begin listening.
	err = d.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	defer func(d *discordgo.Session) {
		err := d.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(d)

	registerCommands(d, command.Sleep, command.Apex, command.Play, command.Hurry, command.Morbius)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
}

func registerCommands(s *discordgo.Session, fs ...command.Factory) {
	var p []string
	for _, f := range fs {
		c, h := f()
		v, err := s.ApplicationCommandCreate(s.State.User.ID, "", c)
		if err != nil {
			log.Errorf("Cannot create '%v' command: %v", c.Name, err)
			continue
		}

		commandHandlers[c.Name] = h
		registeredCommands = append(registeredCommands, v)
		p = append(p, c.Name)
	}
	log.Infof("Started bot with registered commands: %s.", strings.Join(p, ", "))
}
