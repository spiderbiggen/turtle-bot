package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
	"weeb_bot/internal/tenor"
)

type WeightedArgument struct {
	Url      string
	Query    string
	Weight   uint8
	GifCount uint8
	IsSearch bool
}

type Args []*WeightedArgument

func (a Args) Random() *WeightedArgument {
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

type Apex struct{}

func (a *Apex) Command() *discordgo.ApplicationCommand {
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

func (a *Apex) InteractionID() string { return "apex" }

func (a *Apex) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gifCommand("Time for Apex\nLet's go %[1]s\n%[2]s",
		&WeightedArgument{Query: "Apex Legends"},
	)(s, i)
}

type Hurry struct{}

func (h *Hurry) InteractionID() string {
	return "hurry"
}

func (h *Hurry) Command() *discordgo.ApplicationCommand {
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

func (h *Hurry) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gifCommand("Hurry up %[1]s\n%[2]s",
		&WeightedArgument{Query: "hurry up"},
	)(s, i)
}

type Play struct{}

func (p *Play) InteractionID() string {
	return "play"
}

func (p *Play) Command() *discordgo.ApplicationCommand {
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

func (p *Play) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gifCommand("Let's go %[1]s\n%[2]s",
		&WeightedArgument{Query: "games"},
	)(s, i)
}

type Sleep struct{}

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
	gifCommand("%[2]s",
		&WeightedArgument{Query: "sleep", Weight: 80},
		&WeightedArgument{Query: "night", Weight: 70},
		&WeightedArgument{Url: "https://tenor.com/view/frog-dance-animation-cute-funny-gif-17184624", Weight: 1},
	)(s, i)
}

type Morb struct{}

func (c *Morb) InteractionID() string { return "morb" }

func (c *Morb) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "morb",
		Description: "It's morbin time",
	}
}

func (c *Morb) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gifCommand("%[2]s",
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

func gifCommand(gifText string, queries ...*WeightedArgument) Handler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if len(queries) == 0 {
			log.Errorf("No queries for command %s", i.Interaction.ID)
			return
		}
		timeOut := time.NewTimer(10 * time.Second)
		defer timeOut.Stop()

		mention := userFromOptions(s, i).mention()

		c := make(chan string)
		go getGif(c, queries)

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
				s.InteractionResponseDelete(i.Interaction)
				return
			}
			message = fmt.Sprintf(gifText, mention, gif)
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: message})
			if err != nil {
				log.Errorf("discord failed to complete interaction message: %v", err)
				s.InteractionResponseDelete(i.Interaction)
			}
		case <-timeOut.C:
			s.InteractionResponseDelete(i.Interaction)
			log.Warnf("Failed to send gif response within 15 seconds")
		}
	}
}

func getGif(c chan string, queries []*WeightedArgument) {
	q := Args(queries).Random()
	log.Debugf("Using query %+v", q)
	if q.Url != "" {
		c <- q.Url
		return
	}

	for i := 0; i < 5; i++ {
		m := i - 1
		if m > 0 {
			time.Sleep(time.Duration(i*i*25) * time.Millisecond)
		}
		var gifs tenor.ResultList
		var err error
		if q.IsSearch {
			gifs, err = tenor.Search(q.Query, tenor.WithLimit(1))
		} else {
			gifs, err = tenor.Random(q.Query, tenor.WithLimit(1))
		}
		if err != nil {
			log.Errorf("Tenor Failed somewhere. %v", err)
			continue
		}
		if len(gifs) == 0 {
			log.Warnf("No gifs found for query %s", q.Query)
			c <- ""
			return
		}
		c <- gifs[0].URL
		return
	}
}
