package api

import (
	"net/http"
	"os"
	"weeb_bot/commands"

	"github.com/wafer-bw/disgoslash"
	"github.com/wafer-bw/disgoslash/discord"
)

// GuildIDs holds the list of Guild (server) IDs you would like to register
// a slash command to.
var GuildIDs = []string{os.Getenv("GUILD_ID")}

// Global indicates whether or not a slash command should be registered globally
// across all Guilds the bot has access to.
var Global = false

// Credentials holds your Discord Application's secret credentials.
//
// WARNING - Do not track these secrets in version control.
//
// https://discord.com/developers/applications
var Credentials = &discord.Credentials{
	PublicKey: os.Getenv("PUBLIC_KEY"), // Your Discord Application's Public Key
	ClientID:  os.Getenv("CLIENT_ID"),  // Your Discord Application's Client ID
	Token:     os.Getenv("TOKEN"),      // Your Discord Application's Bot's Token
}

// SlashCommandMap is exported for use with the sync package which will
// register the slash command on Discord.
var SlashCommandMap = disgoslash.NewSlashCommandMap(
	commands.CreateSleepCommand(),
)

// Handler is exported for use as a vercel serverless function
// and acts as the entrypoint for slash command requests.
func Handler(w http.ResponseWriter, r *http.Request) {
	handler := &disgoslash.Handler{SlashCommandMap: SlashCommandMap, Creds: Credentials}
	handler.Handle(w, r)
}
