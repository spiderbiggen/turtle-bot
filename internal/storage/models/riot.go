package models

import "weeb_bot/internal/riot"

type RiotAccount struct {
	riot.Summoner
	riot.Region
}
