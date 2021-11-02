package commands

import (
	"fmt"
	"github.com/wafer-bw/disgoslash"
	"github.com/wafer-bw/disgoslash/discord"
	"weeb_bot/core"
	"weeb_bot/lib/random"
	"weeb_bot/lib/tenor"
)

func play(request *discord.InteractionRequest) *discord.InteractionResponse {
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

	gifs, err := tenor.Top("Games")
	if err != nil {
		return tenorError(err)
	}
	gif := gifs[random.Intn(len(gifs))]
	if user != nil {
		return &discord.InteractionResponse{
			Type: discord.InteractionResponseTypeChannelMessageWithSource,
			Data: &discord.InteractionApplicationCommandCallbackData{
				Content: fmt.Sprintf("Let's go <@%s>\n%s", *user, gif.URL),
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
				Content: fmt.Sprintf("Let's go @here\n%s", gif.URL),
			},
		}
	}

}

var playCommand = &discord.ApplicationCommand{
	Name:              "play",
	Description:       "Tag the channel or someone to come play some games",
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

func CreatePlayCommand() disgoslash.SlashCommand {
	return disgoslash.NewSlashCommand(playCommand, play, core.Global, core.GuildIDs)
}