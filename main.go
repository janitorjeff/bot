package main

import (
	"log"
	"os/signal"
	"syscall"
	"os"

	"github.com/janitorjeff/bot/twitch"
	"github.com/janitorjeff/bot/discord"

	"github.com/joho/godotenv"
)

func main() {
	myEnv, err := godotenv.Read("secrets.env")
	if err != nil {
		log.Fatalf("failed to read enviromental variables: %v\n", err)
	}

	readEnvVar := func(name string) string {
		v, ok := myEnv[name]
		if !ok {
			log.Fatalf("no $%s given\n", name)
		}
		return v
	}

	twitchOauth := readEnvVar("JEFF_TWITCH_OAUTH")
	discordToken := readEnvVar("JEFF_DISCORD_TOKEN")

	channels := []string{"janitorjeff"}
	log.Println("Connecting to Twitch IRC")
	go twitch.IRCInit("JanitorJeff", twitchOauth, channels)

	log.Println("Connecting to Discord")
	if err := discord.Init(discordToken); err != nil {
		log.Fatalf("failed to connect to discord: %v\n", err)
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
