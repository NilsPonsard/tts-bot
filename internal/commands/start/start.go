package start

import (
	"github.com/bwmarrin/discordgo"
	htgotts "github.com/hegedustibor/htgo-tts"
	cli "github.com/jawher/mow.cli"
	"github.com/nilsponsard/tts-bot/pkg/verbosity"
)

type fileHandler struct{}

func (*fileHandler) Play(filename string) error {
	verbosity.Info(filename)
	return nil
}

// setup ping command
func Start(job *cli.Cmd) {

	// function to execute

	job.Action = func() {

		speech := htgotts.Speech{Folder: "audio", Language: "en", Handler: &fileHandler{}}

		err := speech.Speak("test")

		if err != nil {
			verbosity.Error(err)
		}

		discord, err := discordgo.New("Bot " + "authentication token")
		if err != nil {
			verbosity.Error(err)
		}

		verbosity.Info(discord)

	}
}
