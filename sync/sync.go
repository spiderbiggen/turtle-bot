package main

import (
	"github.com/wafer-bw/disgoslash"
	"weeb_bot/api/interactions"
	"weeb_bot/core"
)

func main() {
	syncer := &disgoslash.Syncer{
		Creds:           core.Credentials,
		SlashCommandMap: interactions.SlashCommandMap,
		GuildIDs:        core.GuildIDs,
	}
	syncer.Sync()
}
