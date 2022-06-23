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
	return apexCommand, gifCommand("Time for Apex\nLet's go %[1]s\n%[2]s",
		&WeightedArgument{Query: "Apex Legends"},
	)
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
	return hurryCommand, gifCommand("Hurry up %[1]s\n%[2]s",
		&WeightedArgument{Query: "hurry up"},
	)
}

func (u *User) mention() string {
	if u == nil {
		return "@here"
	} else {
		return fmt.Sprintf("<@%s>", u.ID)
	}
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
	return playCommand, gifCommand("Let's go %[1]s\n%[2]s",
		&WeightedArgument{Query: "games"},
	)
}

func Sleep() (*discordgo.ApplicationCommand, Handler) {
	var sleepCommand = &discordgo.ApplicationCommand{
		Name:        "sleep",
		Description: "Gets a random good night gif",
	}
	return sleepCommand, gifCommand("%[2]s",
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
	return sleepCommand, gifCommand("%[2]s",
		&WeightedArgument{Query: "Morbius"},
		&WeightedArgument{Query: "Morbin"},
	)
}

func gifCommand(gifText string, queries ...*WeightedArgument) Handler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if len(queries) == 0 {
			log.Errorf("No queries for command %s", i.Interaction.ID)
			return
		}
		timeOut := time.NewTimer(15 * time.Second)
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
		for {
			select {
			case gif := <-c:
				if gif == "" {
					return
				}
				message = fmt.Sprintf(gifText, mention, gif)
				_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: message})
				if err != nil {
					log.Errorf("discord failed to complete interaction message: %v", err)
				}
				return
			case <-timeOut.C:
				s.InteractionResponseDelete(i.Interaction)
				log.Warnf("Failed to send gif response within 15 seconds")
				return
			}
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

	var gifs tenor.ResultList
	var err error
	if q.IsSearch {
		gifs, err = tenor.Search(q.Query, tenor.WithLimit(1))
	} else {
		gifs, err = tenor.Random(q.Query, tenor.WithLimit(1))
	}
	if err != nil {
		log.Errorf("Tenor Failed somewhere. %v", err)
		c <- ""
		return
	}
	if len(gifs) == 0 {
		log.Warnf("No gifs found for query %s", q.Query)
		c <- ""
		return
	}
	c <- gifs[0].URL
}
