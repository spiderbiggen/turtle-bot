package main

import (
	"github.com/wafer-bw/disgoslash"

	"weeb_bot/api"
	"weeb_bot/core"
)

func main() {
	syncer := &disgoslash.Syncer{
		Creds:           core.Credentials,
		SlashCommandMap: api.SlashCommandMap,
		GuildIDs:        core.GuildIDs,
	}
	syncer.Sync()
}
