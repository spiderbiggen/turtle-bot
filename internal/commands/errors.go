package commands

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func tenorError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	log.Errorln("Tenor Failed somewhere", err)
	err = s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Tenor could not be reached.",
			},
		},
	)
	if err != nil {
		log.Errorf("discord failed to send error message: %v", err)
	}
}
