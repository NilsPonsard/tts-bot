package interactions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nilsponsard/tts-bot/pkg/verbosity"
)

func InitCommands(session *discordgo.Session) {

	var commands = []*discordgo.ApplicationCommand{
		{
			Name:        "setup",
			Description: "setup no-mic channel and its language",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "language",
					Description: "Language identifier",
					Required:    true,
				},
			},
		},
	}
	var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){

		"setup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{
				// Here we need to convert raw interface{} value to wanted type.
				// Also, as you can see, here is used utility functions to convert the value
				// to particular type. Yeah, you can use just switch type,
				// but this is much simpler
				i.ApplicationCommandData().Options[0].StringValue(),
			}
			msgformat :=
				` Now you just learned how to use command options. Take a look to the value of which you've just entered:
				> language: %s
	`
			if len(i.ApplicationCommandData().Options) >= 4 {
				margs = append(margs, i.ApplicationCommandData().Options[3].ChannelValue(nil).ID)
				msgformat += "> channel-option: <#%s>\n"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, we'll discuss them in "responses" part
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	for _, v := range commands {
		var err error = nil

		verbosity.Debug(v.Name)
		_, err = session.ApplicationCommandCreate(session.State.User.ID, "", v)
		if err != nil {
			verbosity.Error("Cannot create command ", v, " : ", err)
		}
	}

}
