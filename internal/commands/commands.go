package commands

import (
	cli "github.com/jawher/mow.cli"
	"github.com/nilsponsard/tts-bot/internal/commands/start"
)

// configure subcommands
func SetupCommands(app *cli.Cli) {
	app.Command("start", "start", start.Start)
}
