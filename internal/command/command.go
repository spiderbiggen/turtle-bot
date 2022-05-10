package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"time"
	"weeb_bot/internal/random"
	"weeb_bot/internal/tenor"
)

type Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Factory func() (*discordgo.ApplicationCommand, Handler)

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
	return apexCommand, gifCommand("Time for Apex\nLet's go %s", "Time for Apex\nLet's go %s\n%s", true, "Apex Legends")
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
	return hurryCommand, gifCommand("Hurry up %s", "Hurry up %s\n%s", true, "hurry up")
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
	return playCommand, gifCommand("Let's go %s", "Let's go %s\n%s", true, "games")
}

func Sleep() (*discordgo.ApplicationCommand, Handler) {
	var sleepCommand = &discordgo.ApplicationCommand{
		Name:        "sleep",
		Description: "Gets a random good night gif",
	}
	return sleepCommand, gifCommand("Good Night!", "%s", false, "night", "sleep")
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

func gifCommand(baseText, gifText string, withUser bool, queries ...string) Handler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			q := queries[random.Intn(len(queries))]
			gifs, err := tenor.Random(q, tenor.WithLimit(50))
			if err != nil {
				tenorError(s, i, err)
			}
			c <- gifs[random.Intn(len(gifs))]
		}()

		var err error
		var sent bool
		var message string
		for {
			select {
			case gif := <-c:
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
			case <-time.After(2500 * time.Millisecond):
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
