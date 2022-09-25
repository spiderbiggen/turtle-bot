package command

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
	"turtle-bot/internal/riot"
	"turtle-bot/internal/stats"
	"turtle-bot/internal/storage/models"
	"turtle-bot/internal/storage/postgres"
)

var (
	ErrNotLinked = errors.New("not linked")
)

type RiotGroup struct {
	Api *riot.Client
	Db  *postgres.Client
}

func (g *RiotGroup) InteractionID() string { return "lol" }

func (g *RiotGroup) Command() *discordgo.ApplicationCommand {
	var regionOptions []*discordgo.ApplicationCommandOptionChoice
	for _, region := range riot.Regions {
		realm, _ := region.Realm()
		regionOptions = append(regionOptions, &discordgo.ApplicationCommandOptionChoice{
			Name:  realm,
			Value: region,
		})
	}

	return &discordgo.ApplicationCommand{
		Name:        "lol",
		Description: "League of Legends related commands",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "link",
				Description: "link your league account to your discord account",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "summoner",
						Description: "Your Summoner name",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "region",
						Description: "The region of your account",
						Required:    true,
						Choices:     regionOptions,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "average",
				Description: "Get average stats from your games.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "region",
						Description: "The region of your account",
						Choices:     regionOptions,
					},
				},
			},
		},
	}
}

func (g *RiotGroup) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "LoL commands are under construction",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		log.Errorf("discord failed to respond with an error message: %v", err)
	}
	return
	// TODO Under construction
	log.Debugf("Responding to lol.%s", options[0].Name)
	switch options[0].Name {
	case "link":
		g.linkHandler(s, i)
	case "average":
		g.averageHandler(s, i)
	default:
		content := "Oops, something went wrong.\n" +
			"Hol' up, you aren't supposed to see this message."
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Errorf("discord failed to respond with an error message: %v", err)
		}
	}
}

func (g *RiotGroup) linkHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var region riot.Region
	var summonerName string
	subOptions := i.Interaction.ApplicationCommandData().Options[0].Options
	for _, option := range subOptions {
		if option.Name == "summoner" {
			summonerName = option.StringValue()
		} else if option.Name == "region" {
			region = riot.Region(option.UintValue())
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Errorf("discord failed to defer interaction: %v", err)
		return
	}
	cs := make(chan string)
	ce := make(chan string)

	go func() {
		summoner, err := g.Api.SummonerByName(ctx, region, summonerName)
		if err != nil {
			log.Errorf("User did not enter correct summonerName or Region, %v", err)
			realm, _ := region.Realm()
			ce <- fmt.Sprintf("No user found for %s in %s", summonerName, realm)
			return
		}
		usr := i.User
		if usr == nil {
			usr = i.Member.User
		}
		acc := models.RiotAccount{
			Summoner: *summoner,
			Region:   region,
		}
		err = g.Db.InsertDiscordSummoner(ctx, usr.ID, i.ChannelID, acc)
		if err != nil {
			log.Errorf("Failed to store summoner in database, %v", err)
			ce <- "Internal server error, please try again later."
			return
		}
		go g.GetMatchHistory(acc)
		cs <- fmt.Sprintf("Successfully linked %s to you", summoner.SummonerName)
	}()

	select {
	case msg := <-cs:
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
		if err != nil {
			_ = s.InteractionResponseDelete(i.Interaction)
			log.Errorf("discord failed to respond with a message: %v", err)
		}
	case msg := <-ce:
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
		if err != nil {
			log.Errorf("discord failed to respond with an error message: %v", err)
		}
	case <-ctx.Done():
		_ = s.InteractionResponseDelete(i.Interaction)
		log.Warnf("Failed to send subscribe response within 15 seconds")
	}
}

func (g *RiotGroup) averageHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})
	if err != nil {
		log.Errorf("discord failed to defer interaction: %v", err)
		return
	}
	type response struct {
		Stats    stats.StatMap
		Summoner riot.Summoner
	}

	cs := make(chan response)
	ce := make(chan error)

	go func() {
		usr := i.User
		if usr == nil {
			usr = i.Member.User
		}

		_, summoner, err := g.Db.GetDiscordSummoner(ctx, usr.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = ErrNotLinked
			}
			ce <- err
			return
		}
		// TODO link to league-stat-server
		cs <- response{Stats: nil, Summoner: summoner.Summoner}
	}()

	select {
	case resp := <-cs:
		name := fmt.Sprintf("Averages %s %s", resp.Summoner.SummonerName, time.Now())
		_ = s.InteractionResponseDelete(i.Interaction)
		thread, err := s.ThreadStart(i.ChannelID, name, discordgo.ChannelTypeGuildPublicThread, 60)
		if err != nil {
			log.Errorf("discord failed to start thread: %v", err)
		}
		builder := strings.Builder{}
		for s2, f := range resp.Stats {
			if f.Min == 0 && f.Max == 0 {
				continue
			}
			var next string
			if f.Min == 0 {
				next = fmt.Sprintf("%s=\tavg: %.2f,\tmax: %.2f\n", s2, f.Average(), f.Max)
			} else {
				next = fmt.Sprintf("%s=\tmin: %.2f,\tavg: %.2f,\tmax: %.2f\n", s2, f.Min, f.Average(), f.Max)
			}
			if builder.Len()+len(next) > 2000 {
				_, _ = s.ChannelMessageSend(thread.ID, builder.String())
				builder.Reset()
			} else {
				builder.WriteString(next)
			}
		}
		_, _ = s.ChannelMessageSend(thread.ID, builder.String())
	case err := <-ce:
		if errors.Is(err, ErrNotLinked) {
			msg := "Please use `/lol link` to link your account and wait a few minutes for your matches to be collected."
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
			if err != nil {
				_ = s.InteractionResponseDelete(i.Interaction)
				log.Errorf("discord failed to respond with a message: %v", err)
			}
		} else {
			_ = s.InteractionResponseDelete(i.Interaction)
			log.Errorf("discord failed to respond with an error message: %v", err)
		}
	case <-ctx.Done():
		_ = s.InteractionResponseDelete(i.Interaction)
		log.Warnf("Failed to send average response within 15 seconds")
	}
}

func (g *RiotGroup) GetMatchHistory(acc models.RiotAccount) {
	// TODO move to league-stat-server
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//start := time.Now().AddDate(0, -3, 0)
	//start = time.Date(start.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	//options := &riot.MatchIdsOptions{StartTime: start.Unix(), Count: 100}
	//var count, failed atomic.Int64
	//defer func() {
	//	log.Debugf("Got Matches for %s: saved %d, failed %d", acc.SummonerName, count.Load(), failed.Load())
	//}()
	//
	//matchQueue := make([]string, 0, 100)
	//for {
	//	matches, err := g.Api.MatchIds(ctx, acc.Region, acc.Puuid, options)
	//	if err != nil {
	//		log.Errorf("Failed to get match ids for %s: %v", acc.SummonerName, err)
	//		return
	//	}
	//	log.Debugf("%d <- %#v", len(matches), options)
	//
	//	// Remove retrieved matches
	//	ids, err := g.Couch.FilterMatchIds(ctx, matches)
	//	if err != nil {
	//		log.Errorf("Failed to get matches for %s: %v", acc.SummonerName, err)
	//		return
	//	}
	//	matchQueue = append(matchQueue, ids...)
	//
	//	if len(matches) < int(options.Count) {
	//		break
	//	}
	//
	//	options.Start += int32(options.Count)
	//}
	//g.Queue.AddMatchIds(matchQueue...)
	//g.Queue.Start()
}
