package main

import (
	"context"
	cRand "crypto/rand"
	"fmt"
	mRand "math/rand"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	rock     = "rock"
	paper    = "paper"
	scissors = "scissors"

	drawStampsPrize = 2
	loseStampsPrize = 1
	winStampsPrize  = 5
)

type gameSeed struct {
	data      [8]byte
	wasSeeded bool
}

var seed *gameSeed

var playMap = map[uint32]string{
	0: rock,
	1: paper,
	2: scissors,
}

var winAssets = map[string][]string{
	rock:     []string{"assets/win_rock_0.png", "assets/win_rock_1.png", "assets/win_rock_2.png"},
	paper:    []string{"assets/win_paper_0.png", "assets/win_paper_1.png", "assets/win_paper_2.png"},
	scissors: []string{"assets/win_scissors_0.png", "assets/win_scissors_1.png", "assets/win_scissors_2.png"},
}

var loseAssets = map[string][]string{
	rock:     []string{"assets/lose_rock_0.png", "assets/lose_rock_1.png", "assets/lose_rock_2.png"},
	paper:    []string{"assets/lose_paper_0.png", "assets/lose_paper_1.png", "assets/lose_paper_2.png"},
	scissors: []string{"assets/lose_scissors_0.png", "assets/lose_scissors_1.png", "assets/lose_scissors_2.png"},
}

var drawAssets = map[string][]string{
	rock:     []string{"assets/draw_rock_0.png", "assets/draw_rock_1.png", "assets/draw_rock_2.png"},
	paper:    []string{"assets/draw_paper_0.png", "assets/draw_paper_1.png", "assets/draw_paper_2.png"},
	scissors: []string{"assets/draw_scissors_0.png", "assets/draw_scissors_1.png", "assets/draw_scissors_2.png"},
}

func SeedNumberGenerator() {
	if seed != nil {
		// TODO: add periodic reseeding
		if seed.wasSeeded == true {
			return
		}
	}

	seed = &gameSeed{
		wasSeeded: false,
	}

	_, err := cRand.Read(seed.data[:])

	if err != nil {
		panic("unable to create seed for rng")
	}

	seed.wasSeeded = true

	return
}

func PlayJanken(s *discordgo.Session, m *discordgo.MessageCreate) {
	ctx := context.Background()
	player := InitPlayer(m.Author)

	if m.Author.ID == s.State.User.ID {
		return
	}

	for i, mention := range m.Mentions {
		if mention == s.State.User {
			break
		} else if i == len(m.Mentions) {
			fmt.Println("exiting early")
			return
		}
	}

	msg := strings.Split(m.Message.Content, " ")

	if len(msg) != 2 {
		return
	}

	attempt := msg[1]
	play := playMap[mRand.Uint32()%3]

	fmt.Printf("Username %s ID %s played %s\n", m.Author.Username, m.Author.ID, attempt)
	fmt.Printf("cpu played %s\n", play)

	if play == attempt {
		sendDrawMessage(ctx, s, m, attempt, player)
		return
	}

	if attempt == rock {
		if play == scissors {
			sendWinMessage(ctx, s, m, attempt, player)
		} else {
			sendLoseMessage(ctx, s, m, attempt, player)
		}

		return
	}

	if attempt == paper {
		if play == rock {
			sendWinMessage(ctx, s, m, attempt, player)
		} else {
			sendLoseMessage(ctx, s, m, attempt, player)
		}

		return
	}

	if attempt == scissors {
		if play == paper {
			sendWinMessage(ctx, s, m, attempt, player)
		} else {
			sendLoseMessage(ctx, s, m, attempt, player)
		}

		return
	}
}

func channelMessageSendComplexWithAsset(s *discordgo.Session, m *discordgo.MessageCreate, assetPath string, f *os.File, resultMessage string) {
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: resultMessage,
		Files: []*discordgo.File{
			&discordgo.File{
				Name:   assetPath,
				Reader: f,
			},
		},
	})
}

// TODO: move function to instead be Player method
func printSingularOrPluralStamp(stampsPrize int64) string {
	if stampsPrize == 1 {
		return "stamp"
	}

	return "stamps"
}

// TODO: move function to instead be Player method
func generateResultMessage(playerMention string, stampsPrize int64, stampsCount int64) string {
	return fmt.Sprintf(
		"%s You've earned %d %s! Your new stamp total is %d %s.",
		playerMention,
		stampsPrize,
		printSingularOrPluralStamp(stampsPrize),
		stampsCount,
		printSingularOrPluralStamp(stampsCount),
	)
}

// TODO: move all send*Message functions to instead be Player methods
// TODO: add log levels like WARN, maybe via zap
func sendWinMessage(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, playerAttempt string, player *Player) {
	player.AddStampsWin(ctx, m.Author)
	resultMessage := generateResultMessage(m.Author.Mention(), winStampsPrize, player.StampsCount)

	assetPath := winAssets[playerAttempt][mRand.Uint32()%3]
	f, err := os.Open(assetPath)

	if err != nil {
		fmt.Println("unable to read image from disk")
		s.ChannelMessageSend(m.ChannelID, resultMessage)
		return
	}

	defer f.Close()

	channelMessageSendComplexWithAsset(s, m, assetPath, f, resultMessage)
}

func sendLoseMessage(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, playerAttempt string, player *Player) {
	player.AddStampsLose(ctx, m.Author)
	resultMessage := generateResultMessage(m.Author.Mention(), loseStampsPrize, player.StampsCount)

	assetPath := loseAssets[playerAttempt][mRand.Uint32()%3]
	f, err := os.Open(assetPath)

	if err != nil {
		fmt.Println("unable to read lose image from disk")
		s.ChannelMessageSend(m.ChannelID, resultMessage)
		return
	}

	defer f.Close()

	channelMessageSendComplexWithAsset(s, m, assetPath, f, resultMessage)
}

func sendDrawMessage(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, playerAttempt string, player *Player) {
	player.AddStampsDraw(ctx, m.Author)
	resultMessage := generateResultMessage(m.Author.Mention(), drawStampsPrize, player.StampsCount)

	assetPath := drawAssets[playerAttempt][mRand.Uint32()%3]
	f, err := os.Open(assetPath)

	if err != nil {
		fmt.Println("unable to read draw image from disk")
		s.ChannelMessageSend(m.ChannelID, resultMessage)
		return
	}

	defer f.Close()

	channelMessageSendComplexWithAsset(s, m, assetPath, f, resultMessage)
}
