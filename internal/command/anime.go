package command

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"time"
	kitsuApi "turtle-bot/internal/kitsu"
	"turtle-bot/internal/storage/postgres"
)

type animeGroup struct {
	*kitsuApi.Client
	db *postgres.Client
}

func AnimeGroup(client *kitsuApi.Client, db *postgres.Client) InteractionHandler {
	return &animeGroup{Client: client, db: db}
}

func (a *animeGroup) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "anime",
		Description: "Anime related commands",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "search",
				Description: "Search anime on kitsu",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "title",
						Description: "Anime title",
						Required:    true,
					},
				},
			},
		},
	}
}

func (a *animeGroup) InteractionID() string {
	return "anime"
}

func (a *animeGroup) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	log.Debugf("Responding to anime.%s", options[0].Name)
	switch options[0].Name {
	case "search":
		a.searchHandler(s, i)
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

func (a *animeGroup) searchHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	subOptions := i.Interaction.ApplicationCommandData().Options[0].Options
	var title string
	for _, option := range subOptions {
		if option.Name == "title" {
			title = option.StringValue()
		}
	}
	if title == "" {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide a title",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Errorf("Failed to respond with error: %v", err)
		}
		return
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

	cr := make(chan []*kitsuApi.Anime)
	ce := make(chan string)

	go func() {
		results, err := a.SearchAnime(ctx, title)
		if err != nil {
			log.Errorf("Failed to search anime: %v", err)
			ce <- fmt.Sprintf("Couldn't find anime for search term %s", title)
			return
		}
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			for _, anime := range results {
				var image string
				if imageUrl := coverImage(anime.Cover, anime.Poster); imageUrl != "" {
					image = imageUrl
				}
				var sTime sql.NullTime
				if anime.CreatedAt != nil {
					sTime.Time = *anime.CreatedAt
					sTime.Valid = true
				}
				err := a.db.InsertAnime(ctx, postgres.Anime{
					ID:             anime.ID,
					CanonicalTitle: anime.CanonicalTitle,
					QueryTitle:     anime.CanonicalTitle,
					ImageURL:       image,
					CreatedAt:      sTime,
				})
				if err != nil {
					log.Warnf("Failed to insert anime result: %v", err)
				}
			}
		}()

		cr <- results
	}()

	select {
	case results := <-cr:
		resultCount := len(results)
		if resultCount > 5 {
			resultCount = 6
		}
		options := make([]discordgo.SelectMenuOption, 0, resultCount)
		for _, result := range results[:resultCount] {
			options = append(options, discordgo.SelectMenuOption{
				Label: result.CanonicalTitle,
				Value: result.ID,
			})
		}
		msg := "Select result below"
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
			Components: &[]discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID: "anime_select",
						Options:  options,
					},
				}},
			},
		})
		if err != nil {
			_ = s.InteractionResponseDelete(i.Interaction)
			log.Errorf("discord failed to respond with a message: %v", err)
		}
	case msg := <-ce:
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
		if err != nil {
			_ = s.InteractionResponseDelete(i.Interaction)
			log.Errorf("discord failed to respond with an error message: %v", err)
		}
	case <-ctx.Done():
		_ = s.InteractionResponseDelete(i.Interaction)
		log.Warnf("Failed to send search response within 15 seconds")
	}
}

func coverImage(i ...*kitsuApi.ImageSet) string {
	for _, imageSet := range i {
		if imageSet == nil {
			continue
		}
		if imageSet.Medium != nil {
			return *imageSet.Medium
		} else {
			return imageSet.Original
		}
	}
	return ""
}

func (a *animeGroup) ComponentID() string {
	return "anime_select"
}

func (a *animeGroup) HandleComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	data := i.MessageComponentData()
	if len(data.Values) == 0 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please select an option",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Errorf("Failed to respond with error: %v", err)
		}
		return
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		_ = s.InteractionResponseDelete(i.Interaction)
		log.Errorf("Failed to defer anime select component: %v", err)
		return
	}

	cs := make(chan struct{})
	ce := make(chan string)

	id := data.Values[0]
	go func() {
		sub := postgres.AnimeSubscription{AnimeID: id, GuildID: i.GuildID, ChannelID: i.ChannelID}
		err := a.db.InsertAnimeSubscription(ctx, sub)
		if err != nil {
			log.Errorf("Failed to insert subscription: %v", err)
			_ = s.InteractionResponseDelete(i.Interaction)
			ce <- fmt.Sprintf("Failed to subscribe to %s\nInternal server error", id)
			return
		}
		cs <- struct{}{}
	}()

	select {
	case <-cs:
		msg := fmt.Sprintf("Successfully subscribed to %s", id)
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		if err != nil {
			_ = s.InteractionResponseDelete(i.Interaction)
			log.Errorf("discord failed to respond with an error message: %v", err)
		}
	case msg := <-ce:
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
		if err != nil {
			_ = s.InteractionResponseDelete(i.Interaction)
			log.Errorf("discord failed to respond with an error message: %v", err)
		}
	case <-ctx.Done():
		_ = s.InteractionResponseDelete(i.Interaction)
		log.Warnf("Failed to send search response within 15 seconds")
	}
}
