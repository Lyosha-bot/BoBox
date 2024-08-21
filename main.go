package main

import (
	"Bobox/game/game_session"
	"Bobox/handlers"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't get .env")
	}

	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal("Couldn't create session")
	}

	initHandlers(session)

	session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentGuildMessages | discordgo.IntentGuildMessageReactions

	err = session.Open()
	if err != nil {
		log.Fatalf("Couldn't open session: %s", err.Error())
	}
	defer session.Close()

	log.Print("Bot is up!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func initHandlers(session *discordgo.Session) {
	game_session.Init(session)
	session.AddHandler(handlers.OnMessageCreate)
	session.AddHandler(handlers.OnMessageReactionAdd)
}
