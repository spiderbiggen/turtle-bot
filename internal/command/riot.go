package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
	"weeb_bot/internal/riot"
	"weeb_bot/internal/storage/postgres"
)

type riotGroup struct {
	api *riot.Client
	db  *postgres.Client
}

func RiotGroup(api *riot.Client, store *postgres.Client) InteractionHandler {
	return &riotGroup{
		api: api,
		db:  store,
	}
}

func (g *riotGroup) InteractionID() string { return "lol" }

func (g *riotGroup) Command() *discordgo.ApplicationCommand {
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

func (g *riotGroup) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

func (g *riotGroup) linkHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		summoner, err := g.api.SummonerByName(ctx, region, summonerName)
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
		err = g.db.InsertDiscordSummoner(ctx, usr.ID, summoner)
		if err != nil {
			log.Errorf("Failed to store summoner in database, %v", err)
			ce <- "Internal server error, please try again later."
			return
		}
		cs <- fmt.Sprintf("Successfully linked %s to you", summoner.SummonerName)
	}()

	select {
	case msg := <-cs:
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: msg})
		if err != nil {
			s.InteractionResponseDelete(i.Interaction)
			log.Errorf("discord failed to respond with a message: %v", err)
		}
	case msg := <-ce:
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: msg})
		if err != nil {
			log.Errorf("discord failed to respond with an error message: %v", err)
		}
	case <-ctx.Done():
		s.InteractionResponseDelete(i.Interaction)
		log.Warnf("Failed to send subscribe response within 15 seconds")
	}
}
