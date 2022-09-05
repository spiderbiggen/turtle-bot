package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
	"turtle-bot/internal/tenor"
)

type WeightedArgument struct {
	Url      string
	Query    string
	Weight   uint8
	GifCount uint8
	IsSearch bool
}

type Args []*WeightedArgument

func (a Args) Pick() *WeightedArgument {
	switch len(a) {
	case 0:
		return nil
	case 1:
		return a[0]
	}

	var sum int
	for _, argument := range a {
		if argument.Weight == 0 {
			argument.Weight = 1
		}
		sum += int(argument.Weight)
	}
	weight := rand.Intn(sum)
	for _, argument := range a {
		weight -= int(argument.Weight)
		if weight < 0 {
			return argument
		}
	}
	return a[len(a)-1]
}

type Apex struct {
	*tenor.Client
	Cache *cache.Cache
}

func (c *Apex) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "apex",
		Description: "Drops an apex gif with someones name",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user you want to summon",
				Required:    false,
			},
		},
	}
}

func (c *Apex) InteractionID() string { return "apex" }

func (c *Apex) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gifCommand(c.Client, c.Cache, "Time for Apex\nLet's go %[1]s\n%[2]s",
		&WeightedArgument{Query: "Apex Legends"},
	)(s, i)
}

type Hurry struct {
	*tenor.Client
	Cache *cache.Cache
}

func (c *Hurry) InteractionID() string {
	return "hurry"
}

func (c *Hurry) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "hurry",
		Description: "Hurry up",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user you want to summon",
				Required:    false,
			},
		},
	}
}

func (c *Hurry) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gifCommand(c.Client, c.Cache, "Hurry up %[1]s\n%[2]s",
		&WeightedArgument{Query: "hurry up"},
	)(s, i)
}

type Play struct {
	*tenor.Client
	Cache *cache.Cache
}

func (c *Play) InteractionID() string {
	return "play"
}

func (c *Play) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "play",
		Description: "Tag the channel or someone to come play some games",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user you want to summon",
				Required:    false,
			},
		},
	}
}

func (c *Play) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gifCommand(c.Client, c.Cache, "Let's go %[1]s\n%[2]s",
		&WeightedArgument{Query: "games"},
	)(s, i)
}

type Sleep struct {
	*tenor.Client
	Cache *cache.Cache
}

func (c *Sleep) InteractionID() string {
	return "sleep"
}

func (c *Sleep) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "sleep",
		Description: "Gets a random good night gif",
	}
}

func (c *Sleep) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gifCommand(c.Client, c.Cache, "%[2]s",
		&WeightedArgument{Query: "sleep", Weight: 80},
		&WeightedArgument{Query: "night", Weight: 70},
		&WeightedArgument{Url: "https://tenor.com/view/frog-dance-animation-cute-funny-gif-17184624", Weight: 1},
	)(s, i)
}

type Morb struct {
	*tenor.Client
	Cache *cache.Cache
}

func (c *Morb) InteractionID() string { return "morb" }

func (c *Morb) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "morb",
		Description: "It's morbin time",
	}
}

func (c *Morb) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gifCommand(c.Client, c.Cache, "%[2]s",
		&WeightedArgument{Query: "Morbius"},
		&WeightedArgument{Query: "Morbin"},
		&WeightedArgument{Query: "Morb"},
	)(s, i)
}

func (u *User) mention() string {
	if u == nil {
		return "@here"
	} else {
		return fmt.Sprintf("<@%s>", u.ID)
	}
}

func gifCommand(tenor *tenor.Client, memCache *cache.Cache, gifText string, queries ...*WeightedArgument) Handler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if len(queries) == 0 {
			log.Errorf("No queries for command %s", i.Interaction.ID)
			return
		}

		mention := userFromOptions(s, i).mention()

		c := make(chan string)
		go getGif(ctx, tenor, memCache, c, queries)

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		if err != nil {
			log.Errorf("discord failed to defer interaction: %v", err)
		}

		var message string
		select {
		case gif := <-c:
			if gif == "" {
				_ = s.InteractionResponseDelete(i.Interaction)
				return
			}
			message = fmt.Sprintf(gifText, mention, gif)
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &message})
			if err != nil {
				log.Errorf("discord failed to complete interaction message: %v", err)
				_ = s.InteractionResponseDelete(i.Interaction)
			}
		case <-ctx.Done():
			_ = s.InteractionResponseDelete(i.Interaction)
			log.Warnf("Failed to send gif response within 15 seconds")
		}
	}
}

func getGif(ctx context.Context, t *tenor.Client, m *cache.Cache, c chan string, queries []*WeightedArgument) {
	q := Args(queries).Pick()
	log.Debugf("Using query %+v", q)
	if q.Url != "" {
		c <- q.Url
		return
	}

	var gifs tenor.ResultList
	var err error
	if p, found := m.Get(q.Query); found {
		if cast, ok := p.(tenor.ResultList); ok {
			gifs = cast
		}
	}

	if len(gifs) == 0 {
		gifs, err = t.Search(ctx, q.Query, tenor.WithLimit(50), tenor.WithRandom(!q.IsSearch))
		if err != nil {
			log.Errorf("failed to search for gifs: %v", err)
			c <- ""
			return
		}
		if len(gifs) > 0 {
			m.Set(q.Query, gifs, cache.DefaultExpiration)
		}
	}
	if len(gifs) == 0 {
		log.Warnf("No gifs found for query %s", q.Query)
		c <- ""
		return
	}
	c <- gifs[rand.Intn(len(gifs))].URL
}
