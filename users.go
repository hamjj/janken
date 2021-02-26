package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fireStore struct {
	client *firestore.Client
}

var jankenFireStore *fireStore

type Player struct {
	ID          string    `firestore:"id"`
	LastPlayed  time.Time `firestore:"lastPlayed"`
	StampsCount int64     `firestore:"stampsCount"`
}

func InitPlayer(user *discordgo.User) *Player{
	return &Player{
		ID: user.ID,
		LastPlayed: time.Now(),
	}
}

func (p *Player) AddStampsDraw(ctx context.Context, player *discordgo.User) error {
	return p.updateOrCreatePlayer(ctx, drawStampsPrize)
}

func (p *Player) AddStampsLose(ctx context.Context, player *discordgo.User) error {
	return p.updateOrCreatePlayer(ctx, loseStampsPrize)
}

func (p *Player) AddStampsWin(ctx context.Context, player *discordgo.User) error {
	// TODO: abstract out into separate testable function
	return p.updateOrCreatePlayer(ctx, winStampsPrize)
}

func (p *Player) updateOrCreatePlayer(ctx context.Context, stampsPrize int64) error {
	players := jankenFireStore.client.Collection("players")
	docRef := players.Doc(p.ID)

	docSnap, err := docRef.Get(ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			// TODO: cache *DocumentRef.ID to speed up future lookups
			p.StampsCount = stampsPrize
			_, err := docRef.Create(ctx, p)

			// TODO: add retry logic
			// TODO: add caching layer
			if err != nil {
				log.Printf("unable to write to firestore: %s", err)
			}

			// log.Printf("doc ref id: %s", docRef.ID)

			return nil
		} 

		log.Printf("error encountered trying to read from firestore: %s", err)
		return err
	}

	docSnap.DataTo(p)
	p.StampsCount += stampsPrize

	docRef.Update(ctx, []firestore.Update{
		{Path: "stampsCount", Value: p.StampsCount},
		{Path: "lastPlayed", Value: time.Now()},
	})

	return nil
}

func InitDataStores(ctx context.Context) {
	projectID := os.Getenv("GCP_PROJECT_ID")

	if jankenFireStore != nil {
		fmt.Println("firestore client already created, shouldn't see this")
		return
	}

	c, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		panic(fmt.Errorf("unable to connect to firestore: %s", err))
	}

	jankenFireStore = &fireStore{
		client: c,
	}

	fmt.Println("firestore client created")
}
