package interactions

import (
	"fmt"
	"net/http"
	"time"
	"weeb_bot/commands"
	"weeb_bot/core"

	"github.com/wafer-bw/disgoslash"
)

// SlashCommandMap is exported for use with the sync package which will
// register the slash command on Discord.
var SlashCommandMap = disgoslash.NewSlashCommandMap(
	commands.CreateSleepCommand(),
	commands.CreateApexCommand(),
	commands.CreatePlayCommand(),
)

// Handler is exported for use as a vercel serverless function
// and acts as the entrypoint for slash command requests.
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%v] %v\n", time.Now(), core.Credentials)
	handler := &disgoslash.Handler{SlashCommandMap: SlashCommandMap, Creds: core.Credentials}
	handler.Handle(w, r)
}
