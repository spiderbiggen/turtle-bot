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

	gifs := tenor.Top("Apex Legends")
	gif := gifs[random.Intn(len(gifs))]
	if user != nil {
		return &discord.InteractionResponse{
			Type: discord.InteractionResponseTypeChannelMessageWithSource,
			Data: &discord.InteractionApplicationCommandCallbackData{
				Content: fmt.Sprintf("Time for Apex\nLet's go <@%s>\n%s", *user, gif.URL),
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
				Content: gif.URL,
			},
		}
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
