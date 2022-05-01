package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"weeb_bot/internal/random"
	"weeb_bot/internal/tenor"
)

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

func apexHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gifs, err := tenor.Search("Apex Legends", tenor.WithLimit(50))
	if err != nil {
		tenorError(s, i, err)
		return
	}
	gif := gifs[random.Intn(len(gifs))]

	var users []string
	mention := "@here"
	if user := UserFromOptions(s, i); user != nil {
		mention = fmt.Sprintf("<@%s>", user.ID)
		users = append(users, user.ID)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Time for Apex\nLet's go %s\n%s", mention, gif.URL),
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
				Users: users,
			},
		},
	})
	if err != nil {
		log.Errorf("discord failed to send response message: %v", err)
	}
}

func CreateApexCommand() (*discordgo.ApplicationCommand, func(*discordgo.Session, *discordgo.InteractionCreate)) {
	return apexCommand, apexHandler
}
