package commands

import (
	"fmt"
	"github.com/wafer-bw/disgoslash"
	"github.com/wafer-bw/disgoslash/discord"
	"weeb_bot/core"
	"weeb_bot/lib/random"
	"weeb_bot/lib/tenor"
)

func aoe(request *discord.InteractionRequest) *discord.InteractionResponse {
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

	gifs, err := tenor.Top("Age of Empires")
	if err != nil {
		return tenorError(err)
	}
	gif := gifs[random.Intn(len(gifs))]
	if user != nil {
		return &discord.InteractionResponse{
			Type: discord.InteractionResponseTypeChannelMessageWithSource,
			Data: &discord.InteractionApplicationCommandCallbackData{
				Content: fmt.Sprintf("Time for AOE\nLet's go <@%s>\n%s", *user, gif.URL),
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

var aoeCommand = &discord.ApplicationCommand{
	Name:              "aoe",
	Description:       "Drops an aoe gif with someones name",
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

func CreateAoeCommand() disgoslash.SlashCommand {
	return disgoslash.NewSlashCommand(aoeCommand, aoe, core.Global, core.GuildIDs)
}
