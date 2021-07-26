package commands

import (
	"github.com/wafer-bw/disgoslash"
	"github.com/wafer-bw/disgoslash/discord"
	"math/rand"
	"weeb_bot/core"
	"weeb_bot/lib/tenor"
)

func sleep(request *discord.InteractionRequest) *discord.InteractionResponse {
	// Your custom code goes here!
	gifs := tenor.Random("Sleep Well")
	gif := gifs[rand.Intn(len(gifs))]
	return &discord.InteractionResponse{
		Type: discord.InteractionResponseTypeChannelMessageWithSource,
		Data: &discord.InteractionApplicationCommandCallbackData{
			Content: gif.URL,
		},
	}
}

var sleepCommand = &discord.ApplicationCommand{
	Name:              "sleep",
	Description:       "Gets a random good night gif",
	DefaultPermission: true,
}

func CreateSleepCommand() disgoslash.SlashCommand {
	return disgoslash.NewSlashCommand(sleepCommand, sleep, core.Global, core.GuildIDs)
}
