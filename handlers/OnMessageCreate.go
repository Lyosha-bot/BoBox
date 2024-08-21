package handlers

import (
	"Bobox/game/game_session"
	"Bobox/lib/embed"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

const (
	UnknownCommandMessage = "Unknown command"
	HelpMessage           = `
BoBox is a small game bot about moving boxes.

Bot prefix is .b.

help - information about bot and commands
play [level_number] - play the game starting from specific level, if no level selected it will begin from the start
stop - to stop current game session
`
	Prefix = ".b"
)

func OnMessageCreate(session *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.Bot {
		return
	}

	commandData := strings.Split(event.Content, " ")
	if len(commandData) == 0 || commandData[0] != Prefix {
		return
	}

	switch commandData[1] {
	case "help":
		session.ChannelMessageSendEmbed(event.ChannelID, embed.Wrap("Help", HelpMessage))
	case "play":
		if len(commandData) == 3 {
			intLevel, err := strconv.Atoi(commandData[2])
			if err == nil {
				game_session.New(event.ChannelID, event.Author.ID, intLevel-1)
			}
		} else {
			game_session.New(event.ChannelID, event.Author.ID, 0)
		}
	case "stop":
		game_session.Stop(event.ChannelID, event.Author.ID)
	default:
		session.ChannelMessageSendEmbed(event.ChannelID, embed.Wrap("Error", UnknownCommandMessage))
	}
}
