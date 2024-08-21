package game_session

import (
	"Bobox/game"
	"Bobox/game/playground"
	"Bobox/lib/e"
	"Bobox/lib/embed"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"strings"
)

const GameExistsMessage = `You already have active game!
Use command "stop" to end previous game`

type Game struct {
	Player     string
	Message    *discordgo.Message
	Level      int
	CanProcess bool
	FieldData  playground.FieldData
}

var BotSession *discordgo.Session
var GamesMap map[string]*Game

func Init(session *discordgo.Session) {
	BotSession = session
	GamesMap = make(map[string]*Game)
}

func New(channelID string, playerID string, level int) {
	_, exists := GamesMap[playerID]
	if exists {
		_, err := BotSession.ChannelMessageSendEmbed(channelID, embed.Wrap("Failed to start the game", GameExistsMessage))
		if err != nil {
			log.Printf("Couldn't send game message for %s", playerID)
		}
		return
	}

	gameMessage, err := BotSession.ChannelMessageSendEmbed(channelID, embed.Wrap("Starting up...", "Generating new game..."))
	if err != nil {
		log.Printf("Couldn't start game for %s", playerID)
		return
	}

	newGame := &Game{
		Player:     playerID,
		Message:    gameMessage,
		Level:      level,
		CanProcess: false,
		FieldData:  playground.New(),
	}

	GamesMap[playerID] = newGame

	newGame.addReactions()

	if err := newGame.LoadLevel(); err != nil {
		fmt.Printf("Couldn't load level: %s", err.Error())
	}
}

func Stop(channelID string, playerID string) {
	_, exists := GamesMap[playerID]
	if !exists {
		BotSession.ChannelMessageSendEmbed(channelID, embed.Wrap("Error", "No game session found"))
		return
	}
	delete(GamesMap, playerID)
	BotSession.ChannelMessageSendEmbed(channelID, embed.Wrap("Success", "Previous game session was stopped"))
}

func (g *Game) LoadLevel() (err error) {
	defer func() { err = e.WrapIfErr("couldn't load level", err) }()

	g.CanProcess = false

	err = g.FieldData.LoadData(g.Level)
	if err != nil {
		return err
	}

	if err = g.renderLevel(); err != nil {
		return err
	}

	g.CanProcess = true

	return nil
}

func (g *Game) ProcessMove(move string) bool {
	if !g.CanProcess {
		return false
	}

	if move == game.EmojiRestart {
		err := g.LoadLevel()
		if err != nil {
			log.Printf("Couldn't restart level: %s", err.Error())
		}
		return false
	}

	dir := moveDirection(move)
	if dir.X+dir.Y == 0 {
		return false
	}

	ok := g.move(g.FieldData.PlayerPosition, dir, game.FieldPlayer)

	if ok {
		g.FieldData.SetPosition(g.FieldData.PlayerPosition, game.FieldEmpty)
		g.FieldData.PlayerPosition.X += dir.X
		g.FieldData.PlayerPosition.Y += dir.Y
		_ = g.renderLevel()
		return !g.FieldData.AnyTargetLeft()
	}
	return false
}

func (g *Game) renderLevel() (err error) {
	defer func() { err = e.WrapIfErr("Couldn't render level", err) }()

	var resMessage strings.Builder

	curTarget := 0
	targetPos := g.FieldData.TargetPositions[curTarget]

	resMessage.WriteString(strings.Repeat(fieldEmoji(game.FieldWall), g.FieldData.Size.X+2)) // Border

	for y := 0; y < g.FieldData.Size.Y; y++ {
		resMessage.WriteString("\n")
		resMessage.WriteString(fieldEmoji(game.FieldWall))
		for x := 0; x < g.FieldData.Size.X; x++ {
			object := g.FieldData.ObjectFromPosition(game.Vector{
				X: x,
				Y: y,
			})
			if targetPos.X == x && targetPos.Y == y {
				if object == game.FieldEmpty {
					resMessage.WriteString(fieldEmoji(game.FieldTarget))
				} else {
					resMessage.WriteString(fieldEmoji(object))
				}
				if curTarget < len(g.FieldData.TargetPositions)-1 {
					curTarget++
					targetPos = g.FieldData.TargetPositions[curTarget]
				}
			} else {
				resMessage.WriteString(fieldEmoji(object))
			}
		}
		resMessage.WriteString(fieldEmoji(game.FieldWall))
	}

	resMessage.WriteString("\n")
	resMessage.WriteString(strings.Repeat(fieldEmoji(game.FieldWall), g.FieldData.Size.X+2)) // Border

	_, err = BotSession.ChannelMessageEditEmbed(g.Message.ChannelID, g.Message.ID, embed.Wrap("Level "+strconv.Itoa(g.Level+1), resMessage.String()))
	return err
}

func (g *Game) addReactions() {
	_ = g.addSingleReaction(game.EmojiLeft)
	_ = g.addSingleReaction(game.EmojiUp)
	_ = g.addSingleReaction(game.EmojiDown)
	_ = g.addSingleReaction(game.EmojiRight)
	_ = g.addSingleReaction(game.EmojiRestart)
}

func (g *Game) addSingleReaction(reaction string) error {
	return BotSession.MessageReactionAdd(g.Message.ChannelID, g.Message.ID, reaction)
}

func (g *Game) move(position game.Vector, dir game.Vector, objectID int) bool {
	nextPosition := game.Vector{
		X: dir.X + position.X,
		Y: dir.Y + position.Y,
	}

	if !g.FieldData.IsValidPosition(nextPosition) {
		return false
	}

	nextObjectID := g.FieldData.ObjectFromPosition(nextPosition)

	if nextObjectID == game.FieldEmpty {
		g.FieldData.SetPosition(nextPosition, objectID)
		return true
	} else if nextObjectID == game.FieldBox {
		moveRes := g.move(nextPosition, dir, game.FieldBox)
		if moveRes {
			g.FieldData.SetPosition(nextPosition, objectID)
		}
		return moveRes
	}

	return false
}

func moveDirection(direction string) game.Vector {
	switch direction {
	case game.EmojiUp:
		return game.Vector{X: 0, Y: -1}
	case game.EmojiDown:
		return game.Vector{X: 0, Y: 1}
	case game.EmojiRight:
		return game.Vector{X: 1, Y: 0}
	case game.EmojiLeft:
		return game.Vector{X: -1, Y: 0}
	default:
		return game.Vector{}
	}
}

func fieldEmoji(id int) string {
	switch id {
	case 1:
		return game.EmojiWall
	case 2:
		return game.EmojiPlayer
	case 3:
		return game.EmojiBox
	case 4:
		return game.EmojiTarget
	default:
		return game.EmojiEmpty
	}
}
