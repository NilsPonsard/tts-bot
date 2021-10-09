package start

import (
	"github.com/bwmarrin/discordgo"
	cli "github.com/jawher/mow.cli"
	"github.com/nilsponsard/tts-bot/pkg/verbosity"
)

// setup ping command
func Start(job *cli.Cmd) {

	// function to execute

	job.Action = func() {
		discord, err := discordgo.New("Bot " + "authentication token")
		if err != nil {
			verbosity.Error(err)
		}

		verbosity.Info(discord)

	}
}
