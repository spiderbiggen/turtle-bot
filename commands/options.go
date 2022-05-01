package commands

import (
	"github.com/bwmarrin/discordgo"
)

// UserFromOptions looks for the first option with the "user" key unless another key is provided.
func UserFromOptions(s *discordgo.Session, i *discordgo.InteractionCreate, keys ...string) *discordgo.User {
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
