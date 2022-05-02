package command

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"weeb_bot/internal/random"
	"weeb_bot/internal/tenor"
)

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

func playHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gifs, err := tenor.Search("Games", tenor.WithLimit(50))
	if err != nil {
		tenorError(s, i, err)
		return
	}
	gif := gifs[random.Intn(len(gifs))]

	mention := "@here"
	if user := userFromOptions(s, i); user != nil {
		mention = fmt.Sprintf("<@%s>", user.ID)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Let's go %s\n%s", mention, gif.URL),
		},
	})
	if err != nil {
		log.Errorf("discord failed to send response message: %v", err)
	}
}

func Play() (*discordgo.ApplicationCommand, Handler) {
	return playCommand, playHandler
}
