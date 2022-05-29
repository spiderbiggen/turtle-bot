package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
	"weeb_bot/internal/tenor"
)

type Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Factory func() (*discordgo.ApplicationCommand, Handler)

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

func Apex() (*discordgo.ApplicationCommand, Handler) {
	var apexCommand = &discordgo.ApplicationCommand{
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
	return apexCommand, gifCommand("Time for Apex\nLet's go %s", "Time for Apex\nLet's go %s\n%s", true, &WeightedArgument{Query: "Apex Legends"})
}

func Hurry() (*discordgo.ApplicationCommand, Handler) {
	var hurryCommand = &discordgo.ApplicationCommand{
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
	return hurryCommand, gifCommand("Hurry up %s", "Hurry up %s\n%s", true, &WeightedArgument{Query: "hurry up"})
}

func Play() (*discordgo.ApplicationCommand, Handler) {
	var playCommand = &discordgo.ApplicationCommand{
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
	return playCommand, gifCommand("Let's go %s", "Let's go %s\n%s", true, &WeightedArgument{Query: "games"})
}

func Sleep() (*discordgo.ApplicationCommand, Handler) {
	var sleepCommand = &discordgo.ApplicationCommand{
		Name:        "sleep",
		Description: "Gets a random good night gif",
	}
	return sleepCommand, gifCommand(
		"Good Night!", "%s", false,
		&WeightedArgument{Query: "sleep", Weight: 80},
		&WeightedArgument{Query: "night", Weight: 70},
		&WeightedArgument{Query: "froggers", Weight: 1, GifCount: 1, IsSearch: true},
	)
}

func Morbin() (*discordgo.ApplicationCommand, Handler) {
	var sleepCommand = &discordgo.ApplicationCommand{
		Name:        "morbin",
		Description: "Morbin",
	}
	_, c := Morbius()
	return sleepCommand, c
}

func Morbius() (*discordgo.ApplicationCommand, Handler) {
	var sleepCommand = &discordgo.ApplicationCommand{
		Name:        "morbius",
		Description: "Morbius",
	}
	return sleepCommand, gifCommand(
		"You got morbed", "%s", false,
		&WeightedArgument{Query: "Morbius"},
		&WeightedArgument{Query: "Morbin"},
	)
}

// userFromOptions looks for the first option with the "user" key unless another key is provided.
func userFromOptions(s *discordgo.Session, i *discordgo.InteractionCreate, keys ...string) *discordgo.User {
	key := "user"
	if len(keys) > 0 {
		key = keys[0]
	}

	for _, option := range i.Interaction.ApplicationCommandData().Options {
		if option.Name == key {
			if user := option.UserValue(s); user != nil {
				return user
			}
		}
	}
	return nil
}

func tenorError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	log.Errorf("Tenor Failed somewhere. %v", err)
	err = s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   1 << 6,
				Content: "Tenor could not be reached.",
			},
		},
	)
	if err != nil {
		log.Errorf("discord failed to send error message: %v", err)
	}
}

func gifCommand(baseText, gifText string, withUser bool, queries ...*WeightedArgument) Handler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if len(queries) == 0 {
			log.Errorf("No queries for command %s", i.Interaction.ID)
			return
		}
		tt := time.After(2 * time.Second)

		var mention string
		if withUser {
			if user := userFromOptions(s, i); user != nil {
				mention = fmt.Sprintf("<@%s>", user.ID)
			} else {
				mention = "@here"
			}
		}

		c := make(chan *tenor.Result)
		go func() {
			q := Args(queries).Random()
			log.Debugf("Using query %+v", q)
			var gifs tenor.ResultList
			var err error
			if q.IsSearch {
				gifs, err = tenor.Search(q.Query, tenor.WithLimit(1))
			} else {
				gifs, err = tenor.Random(q.Query, tenor.WithLimit(1))
			}
			if err != nil {
				tenorError(s, i, err)
				c <- nil
				return
			}
			if len(gifs) == 0 {
				log.Warnf("No gifs found for query %s", q.Query)
				c <- nil
				return
			}
			c <- gifs[0]
		}()

		var err error
		var sent bool
		var message string
		for {
			select {
			case gif := <-c:
				if gif == nil {
					return
				}
				if withUser {
					message = fmt.Sprintf(gifText, mention, gif.URL)
				} else {
					message = fmt.Sprintf(gifText, gif.URL)
				}
				if sent {
					_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: message})
				} else {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: message,
						},
					})
				}
				if err != nil {
					log.Errorf("discord failed to complete interaction message: %v", err)
				}
				return
			case <-tt:
				sent = true
				if withUser {
					message = fmt.Sprintf(baseText, mention)
				} else {
					message = baseText
				}
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: message,
					},
				})
				if err != nil {
					log.Errorf("discord failed to send interaction message: %v", err)
				}
			}
		}
	}
}
