# Janken
Janken is a Discord bot written in go that allows users on your server play "rock paper scissors" against the bot. The bot uses a random number generator to determine whether the user won or lost the round. Points are awarded based on whether the user won, lost, or forced a draw. The total point value is saved, so the user can come back and play again to increase their point total.

## Requirements
- go 1.14+
- Discord API token
- Google Cloud project with the Firestore API enabled

## Usage
```
export GCP_PROJECT_ID=$YOUR_GCP_PROJECT_ID
export DISCORD_TOKEN=$YOUR_DISCORD_API_TOKEN
export GOOGLE_APPLICATION_CREDENTIALS=$PATH_TO_GCP_CREDS 

# build and run the go binary
go mod vendor
go build
./janken

# alternatively, use docker
docker build -t janken:latest .
docker run -it --rm \
    -v $PATH_TO_GCP_CREDS \
    -e GCP_PROJECT_ID=$YOUR_GCP_PROJECT_ID \
    -e DISCORD_TOKEN=$YOUR_DISCORD_API_TOKEN \
    -e GOOGLE_APPLICATION_CREDENTIALS=$PATH_TO_GCP_CREDS \
    janken:latest
```

## How to play
@janken [rock|paper|scissors]

## Notes/FAQ
- The bot tries to post win/lose/draw images depending on the outcome of each game. It looks for these images in the `assets/` directory. No images are provided out of the box, but it's possible to add your own by saving your images using the file name scheme found at the top of [game.go](game.go).
- Why is the project titled "Janken"?
  - The project is a chat bot implementation of Konami's [いちかのBEMANI超じゃんけん大会2020](https://p.eagate.573.jp/game/bemani/bjm2020/) (Ichika no BEMANI chou janken taikai 2020) arcade rhythm game event.
  - "Throughout Japanese history there are frequent references to sansukumi-ken, meaning ken (fist) games where "the three who are afraid of one another" (i.e. A beats B, B beats C, and C beats A). This type of game originated in China before being imported to Japan and subsequently also becoming popular among the Japanese. ... Today, the best-known sansukumi-ken is called jan-ken (じゃんけん), which is a variation of the Chinese games introduced in the 17th century. Jan-ken uses the rock, paper, and scissors signs and is the game that the modern version of rock paper scissors derives from directly." [(source)](https://en.wikipedia.org/wiki/Rock_paper_scissors#Origins)