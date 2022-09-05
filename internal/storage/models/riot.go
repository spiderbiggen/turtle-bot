package models

import (
	"turtle-bot/internal/riot"
)

type RiotAccount struct {
	riot.Summoner
	riot.Region
}
