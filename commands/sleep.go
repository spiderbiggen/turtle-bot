package commands

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"weeb_bot/internal/random"
	"weeb_bot/internal/tenor"
)

var queries = [...]string{"night", "sleep"}

var sleepCommand = &discordgo.ApplicationCommand{
	Name:        "sleep",
	Description: "Gets a random good night gif",
}

func sleepHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	q := queries[random.Random().Intn(len(queries))]
	gifs, err := tenor.Random(q, tenor.WithLimit(50))
	if err != nil {
		tenorError(s, i, err)
		return
	}
	gif := gifs[random.Intn(len(gifs))]
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: gif.URL,
		},
	})
	if err != nil {
		log.Errorf("discord failed to send response message: %v", err)
	}
}

func CreateSleepCommand() (*discordgo.ApplicationCommand, func(*discordgo.Session, *discordgo.InteractionCreate)) {
	return sleepCommand, sleepHandler
}
