package worker

import (
	"context"
	"github.com/bwmarrin/discordgo"
)

type Worker func(context.Context, *discordgo.Session)
