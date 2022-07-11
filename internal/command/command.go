package command

import "github.com/bwmarrin/discordgo"

type Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type InteractionHandler interface {
	InteractionID() string
	Command() *discordgo.ApplicationCommand
	HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type ComponentProvider interface {
	Components() []ComponentHandler
}

type ComponentHandler interface {
	ComponentID() string
	HandleComponent(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type User discordgo.User

// userFromOptions looks for the first option with the "user" key unless another key is provided.
func userFromOptions(s *discordgo.Session, i *discordgo.InteractionCreate, keys ...string) *User {
	key := "user"
	if len(keys) > 0 {
		key = keys[0]
	}

	for _, option := range i.Interaction.ApplicationCommandData().Options {
		if option.Name == key {
			if user := option.UserValue(s); user != nil {
				return (*User)(user)
			}
		}
	}
	return nil
}
