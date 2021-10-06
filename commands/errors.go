package commands

import (
	"github.com/wafer-bw/disgoslash/discord"
	"weeb_bot/lib/log"
)

func tenorError(err error) *discord.InteractionResponse {
	log.ErrorLogger.Println("Tenor Failed somewhere")
	return &discord.InteractionResponse{
		Type: discord.InteractionResponseTypeChannelMessageWithSource,
		Data: &discord.InteractionApplicationCommandCallbackData{
			Content: "Tenor could not be reached.",
		},
	}
}
