package start

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"

	"github.com/bwmarrin/discordgo"
	cli "github.com/jawher/mow.cli"
	"github.com/nilsponsard/tts-bot/internal/dgvoice"
	"github.com/nilsponsard/tts-bot/internal/htgotts"
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
		discord.AddHandler(messageCreate)

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

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	verbosity.Debug(m.ChannelID)

	guild, err := s.State.Guild(m.GuildID)

	if err != nil {
		verbosity.Error(err)
		return
	}

	msgChannel, err := s.State.Channel(m.ChannelID)

	if err != nil {
		verbosity.Error(err)
		return
	}
	if msgChannel.Name != "no-mic" {
		return
	}

	var channel *discordgo.Channel

	for _, state := range guild.VoiceStates {

		if state.UserID == m.Author.ID {
			channel, err = s.State.Channel(state.ChannelID)
			if err != nil {
				verbosity.Error(err)
			}
			break
		}

	}

	//-- todo : spak queue

	if channel == nil {
		s.ChannelMessageSend(m.ChannelID, "you are not in a voice channel")
		s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
	} else {

		speech := htgotts.Speech{Folder: "audio", Language: "es"}
		verbosity.Debug(m.Member.Nick)
		file, err := speech.Speak(m.Member.Nick + " a dit : " + m.Content)

		if err != nil {
			verbosity.Error(err)
			return
		}

		ffmpeg := exec.Command("ffmpeg", "-i", file, "-f", "s16le", "-ar", strconv.Itoa(48000), "-ac",
			strconv.Itoa(2), "pipe:1")

		voice, err := s.ChannelVoiceJoin(m.GuildID, channel.ID, false, true)

		if err != nil {
			verbosity.Error(err)
		} else {
			verbosity.Debug(voice)
			voice.Speaking(true)
		}

		out, err := ffmpeg.StdoutPipe()
		if err != nil {
			verbosity.Error(err)
			return
		}
		buffer := bufio.NewReaderSize(out, 16384)
		err = ffmpeg.Start()
		if err != nil {
			verbosity.Error(err)
			return
		}
		voice.Speaking(true)

		soundChan := make(chan []int16, 2)

		go dgvoice.SendPCM(voice, soundChan)
		for {

			audioBuffer := make([]int16, 960*2)
			err = binary.Read(buffer, binary.LittleEndian, &audioBuffer)

			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}

			if err != nil {
				verbosity.Error(err)
				break
			}
			soundChan <- audioBuffer
		}
		voice.Speaking(false)
		s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	}

}
