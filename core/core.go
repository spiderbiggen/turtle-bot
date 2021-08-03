package core

import (
	"github.com/wafer-bw/disgoslash/discord"
	"os"
)

// GuildIDs holds the list of Guild (server) IDs you would like to register
// a slash command to.
var GuildIDs = []string{os.Getenv("GUILD_ID")}

// Global indicates whether a slash command should be registered globally
// across all Guilds the bot has access to.
var Global = os.Getenv("ENV") == "PRODUCTION"

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
