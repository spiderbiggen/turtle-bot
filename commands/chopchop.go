package commands

import (
	"fmt"
	"github.com/wafer-bw/disgoslash"
	"github.com/wafer-bw/disgoslash/discord"
	"weeb_bot/core"
	"weeb_bot/lib/random"
	"weeb_bot/lib/tenor"
)

func chopChop(request *discord.InteractionRequest) *discord.InteractionResponse {
	options := request.Data.Options
	var user *string
	for _, option := range options {
		if option.Name == "name" {
			if _user, ok := option.UserIDValue(); ok {
				user = &_user
			}
			break
		}
	}

	gifs, err := tenor.Top("hurry up")
	if err != nil {
		return tenorError(err)
	}
	gif := gifs[random.Intn(len(gifs))]
	if user != nil {
		return &discord.InteractionResponse{
			Type: discord.InteractionResponseTypeChannelMessageWithSource,
			Data: &discord.InteractionApplicationCommandCallbackData{
				Content: fmt.Sprintf("Hurry up <@%s>\n%s", *user, gif.URL),
				AllowedMentions: &discord.AllowedMentions{
					Parse: []discord.AllowedMentionType{discord.AllowedMentionTypeUserMentions},
					Users: []string{},
				},
			},
		}
	} else {
		return &discord.InteractionResponse{
			Type: discord.InteractionResponseTypeChannelMessageWithSource,
			Data: &discord.InteractionApplicationCommandCallbackData{
				Content: fmt.Sprintf("Hurry up @here\n%s", gif.URL),
			},
		}
	}

}

var chopChopCommand = &discord.ApplicationCommand{
	Name:              "hurry",
	Description:       "Hurry up",
	DefaultPermission: true,
	Options: []*discord.ApplicationCommandOption{
		{
			Type:        discord.ApplicationCommandOptionTypeUser,
			Name:        "name",
			Description: "Enter the name of a user",
			Required:    false,
		},
	},
}

func CreateChopChopCommand() disgoslash.SlashCommand {
	return disgoslash.NewSlashCommand(chopChopCommand, chopChop, core.Global, core.GuildIDs)
}
