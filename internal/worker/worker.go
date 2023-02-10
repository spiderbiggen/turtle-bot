package worker

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

type Worker interface {
	Schedule(*cron.Cron, *discordgo.Session) error
	Run(context.Context, *discordgo.Session) error
}
