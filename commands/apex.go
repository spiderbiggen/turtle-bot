package commands

import (
	"fmt"
	"github.com/wafer-bw/disgoslash"
	"github.com/wafer-bw/disgoslash/discord"
	"weeb_bot/core"
	"weeb_bot/lib/random"
	"weeb_bot/lib/tenor"
)

func apex(request *discord.InteractionRequest) *discord.InteractionResponse {
	gifs, err := tenor.Search("Apex Legends", tenor.WithLimit(50))
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
			Content: fmt.Sprintf("Time for Apex\nLet's go %s\n%s", mention, gif.URL),
			AllowedMentions: &discord.AllowedMentions{
				Parse: []discord.AllowedMentionType{discord.AllowedMentionTypeUserMentions},
				Users: []string{},
			},
		},
	}
}

var apexCommand = &discord.ApplicationCommand{
	Name:              "apex",
	Description:       "Drops an apex gif with someones name",
	DefaultPermission: true,
	Options: []*discord.ApplicationCommandOption{
		{
			Type:        discord.ApplicationCommandOptionTypeUser,
			Name:        "name",
			Description: "Enter the name of the user you want to summon",
			Required:    false,
		},
	},
}

func CreateApexCommand() disgoslash.SlashCommand {
	return disgoslash.NewSlashCommand(apexCommand, apex, core.Global, core.GuildIDs)
}
