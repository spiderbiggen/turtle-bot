package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"sync"
	"sync/atomic"
	"time"
	"weeb_bot/internal/riot"
	"weeb_bot/internal/storage/couch"
	"weeb_bot/internal/storage/postgres"
)

type RiotGroup struct {
	Api   *riot.Client
	Db    *postgres.Client
	Couch *couch.Client
}

func NewRiotGroup(api *riot.Client, store *postgres.Client, couch *couch.Client) InteractionHandler {
	return &RiotGroup{
		Api:   api,
		Db:    store,
		Couch: couch,
	}
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
		},
	}
}

func (g *RiotGroup) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	log.Debugf("Responding to lol.%s", options[0].Name)
	switch options[0].Name {
	case "link":
		g.linkHandler(s, i)
	default:
		content := "Oops, something went wrong.\n" +
			"Hol' up, you aren't supposed to see this message."
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
				Flags:   uint64(discordgo.MessageFlagsEphemeral),
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
			Flags: uint64(discordgo.MessageFlagsEphemeral),
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
		err = g.Db.InsertDiscordSummoner(ctx, usr.ID, summoner)
		if err != nil {
			log.Errorf("Failed to store summoner in database, %v", err)
			ce <- "Internal server error, please try again later."
			return
		}
		go g.GetMatchHistory(summoner, region)
		cs <- fmt.Sprintf("Successfully linked %s to you", summoner.SummonerName)
	}()

	select {
	case msg := <-cs:
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: msg})
		if err != nil {
			_ = s.InteractionResponseDelete(i.Interaction)
			log.Errorf("discord failed to respond with a message: %v", err)
		}
	case msg := <-ce:
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: msg})
		if err != nil {
			log.Errorf("discord failed to respond with an error message: %v", err)
		}
	case <-ctx.Done():
		_ = s.InteractionResponseDelete(i.Interaction)
		log.Warnf("Failed to send subscribe response within 15 seconds")
	}
}

func (g *RiotGroup) GetMatchHistory(s *riot.Summoner, region riot.Region) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	start := time.Now().AddDate(0, -3, 0)
	start = time.Date(start.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	options := &riot.MatchIdsOptions{StartTime: start.Unix(), Count: 100}
	var count, failed atomic.Int64
	defer func() {
		log.Debugf("Got Matches for %s: saved %d, failed %d", s.SummonerName, count.Load(), failed.Load())
	}()

	matchQueue := make([]string, 0, 100)
	for {
		matches, err := g.Api.MatchIds(ctx, region, s.Puuid, options)
		if err != nil {
			log.Errorf("Failed to get match ids for %s: %v", s.SummonerName, err)
			return
		}
		log.Debugf("%d <- %#v", len(matches), options)

		// Remove retrieved matches
		ids, err := g.Couch.FilterMatchIds(ctx, matches)
		if err != nil {
			log.Errorf("Failed to get matches for %s: %v", s.SummonerName, err)
			return
		}
		matchQueue = append(matchQueue, ids...)

		if len(matches) < int(options.Count) {
			break
		}

		options.Start += int32(options.Count)
	}
	wg := sync.WaitGroup{}
	for _, id := range matchQueue {
		go func(id string) {
			wg.Add(1)
			defer wg.Done()
			match, err := g.Api.Match(ctx, region, id)
			if err != nil {
				log.Errorf("Failed to get match %s for %s: %v", id, s.SummonerName, err)
				failed.Add(1)
				return
			}
			log.Debugf("Got match %d", count.Load())
			err = g.Couch.AddMatch(ctx, match)
			if err != nil {
				log.Errorf("Failed to store match %s for %s: %v", id, s.SummonerName, err)
				failed.Add(1)
				return
			}
			count.Add(1)
		}(id)
	}
	wg.Wait()
}
