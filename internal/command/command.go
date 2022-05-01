package command

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Factory func() (*discordgo.ApplicationCommand, Handler)

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
	log.Errorln("Tenor Failed somewhere", err)
	err = s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Tenor could not be reached.",
			},
		},
	)
	if err != nil {
		log.Errorf("discord failed to send error message: %v", err)
	}
}
