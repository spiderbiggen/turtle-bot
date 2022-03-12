package commands

import (
	"github.com/wafer-bw/disgoslash/discord"
)

// UserFromOptions looks for the first option with the "name" key unless another key is provided.
func UserFromOptions(options []*discord.ApplicationCommandInteractionDataOption, keys ...string) *string {
	key := "name"
	if len(keys) > 0 {
		key = keys[0]
	}

	for _, option := range options {
		if option.Name == key {
			if _user, ok := option.UserIDValue(); ok {
				return &_user
			}
		}
	}
	return nil
}
