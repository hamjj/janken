package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)
func main() {
	ctx := context.Background()
	SeedNumberGenerator()
	InitDataStores(ctx)

	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	if err != nil {
		panic("unable to create discord client")
	}

	discord.AddHandler(PlayJanken)

	err = discord.Open()

	if err != nil {
		panic("error opening connection to discord")
	}

	fmt.Println("janken bot is active")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}
