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
	gifs, err := tenor.Search("hurry up", tenor.WithLimit(15))
	if err != nil {
		return tenorError(err)
	}
	gif := gifs[random.Intn(len(gifs))]

	mention := "@here"
	if user := UserFromOptions(request.Data.Options); user != nil {
		mention = fmt.Sprintf("<@%s>", *user)
	}

	return &discord.InteractionResponse{
		Type: discord.InteractionResponseTypeChannelMessageWithSource,
		Data: &discord.InteractionApplicationCommandCallbackData{
			Content: fmt.Sprintf("Hurry up %s\n%s", mention, gif.URL),
			AllowedMentions: &discord.AllowedMentions{
				Parse: []discord.AllowedMentionType{discord.AllowedMentionTypeUserMentions},
			},
		},
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
