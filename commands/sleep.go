package commands

import (
	"github.com/wafer-bw/disgoslash"
	"github.com/wafer-bw/disgoslash/discord"
	"weeb_bot/core"
	"weeb_bot/lib/random"
	"weeb_bot/lib/tenor"
)

var queries = [...]string{"night", "sleep"}

func sleep(_ *discord.InteractionRequest) *discord.InteractionResponse {
	rnd := random.Random()
	q := queries[rnd.Intn(len(queries))]
	gifs, err := tenor.Random(q, tenor.WithLimit(50))
	if err != nil {
		return tenorError(err)
	}
	gif := gifs[random.Intn(len(gifs))]
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
