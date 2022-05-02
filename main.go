package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"weeb_bot/internal/command"
	"weeb_bot/internal/nyaa"
)

var (
	commandHandlers    = make(map[string]command.Handler)
	registeredCommands []*discordgo.ApplicationCommand
)

func main() {
	log.SetReportCaller(true)
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.TraceLevel)
	hostname, _ := os.Hostname()
	log.Warnf("Starting Weeb Bot on %s", hostname)

	d, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("TOKEN")))
	if err != nil {
		log.Fatal(err)
	}
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

	registerCommands(d, command.Sleep, command.Apex, command.Play, command.Hurry)

	//_, err = nyaa.GetAnime(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//c := map[string]string{
	//	"825808364649971712": "825808364649971715",
	//}
	//a := map[string][]string{
	//	"Tate no Yuusha no Nariagari S2": {"825808364649971712"},
	//}
	//for _, group := range groups {
	//	embed := makeEmbed(group)
	//	if guilds, ok := a[group.AnimeTitle]; ok {
	//		for _, guild := range guilds {
	//			if channel, ok := c[guild]; ok {
	//				_, err := d.ChannelMessageSendEmbed(channel, embed)
	//				if err != nil {
	//					log.Fatalln(err)
	//				}
	//			}
	//		}
	//	}
	//}
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

func makeEmbed(g nyaa.Group) *discordgo.MessageEmbed {
	fields := make([]*discordgo.MessageEmbedField, 0, len(g.Downloads))
	for _, d := range g.Downloads {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   d.Resolution,
			Value:  fmt.Sprintf("[torrent](%s)\n[comments](%s)", d.Torrent, d.Comments),
			Inline: true,
		})
	}
	title := g.AnimeTitle
	if g.Episode.Number != 0 {
		title = fmt.Sprintf("%s Ep %d", g.AnimeTitle, g.Episode.Number)
	}
	return &discordgo.MessageEmbed{
		Type:   discordgo.EmbedTypeRich,
		Title:  title,
		Fields: fields,
	}
}
