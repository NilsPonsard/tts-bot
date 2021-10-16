package start

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	cli "github.com/jawher/mow.cli"
	"github.com/nilsponsard/tts-bot/internal/interactions"
	"github.com/nilsponsard/tts-bot/pkg/verbosity"
)

// setup ping command
func Start(job *cli.Cmd) {

	token := job.StringArg("TOKEN", "", "Discord token")

	// function to execute

	job.Action = func() {

		discord, err := discordgo.New("Bot " + *token)
		if err != nil {
			verbosity.Error(err)
			cli.Exit(1)
		}

		verbosity.Info(discord)

		discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
			verbosity.Info("Bot is up!")
		})

		err = discord.Open()
		if err != nil {
			verbosity.Error(err)
			cli.Exit(1)
		}

		interactions.InitCommands(discord)

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		<-stop
		verbosity.Debug("Gracefully shutdowning")

	}
}
