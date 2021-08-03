package main

import (
	"github.com/wafer-bw/disgoslash"
	"os"

	"weeb_bot/api"
	"weeb_bot/core"
)

func main() {
	ids := core.GuildIDs
	if os.Getenv("ENV") == "PRODUCTION" {
		ids = []string{}
	}
	syncer := &disgoslash.Syncer{
		Creds:           core.Credentials,
		SlashCommandMap: api.SlashCommandMap,
		GuildIDs:        ids,
	}
	syncer.Sync()
}
