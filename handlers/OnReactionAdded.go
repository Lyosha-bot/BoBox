package handlers

import (
	"Bobox/game/game_session"
	"github.com/bwmarrin/discordgo"
	"log"
)

func OnMessageReactionAdd(session *discordgo.Session, event *discordgo.MessageReactionAdd) {
	if event.UserID == session.State.User.ID {
		return
	}

	g, ok := game_session.GamesMap[event.UserID]
	if !ok || g.Message.ID != event.MessageID {
		return
	}

	err := session.MessageReactionRemove(event.ChannelID, event.MessageID, event.Emoji.Name, event.UserID)
	if err != nil {
		log.Printf("Couldn't delete reaction %s", event.Emoji.Name)
		return
	}

	if g.ProcessMove(event.Emoji.Name) {
		g.Level++
		g.LoadLevel()
	}
}
