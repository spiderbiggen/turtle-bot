package worker

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
)

type Worker interface {
	Job() *gocron.Job
	Schedule(*gocron.Scheduler, *discordgo.Session) error
	Run(context.Context, *discordgo.Session) error
}
