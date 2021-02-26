The bot tries to post win/lose/draw images depending on the outcome of each game. It looks for these images in the `assets/` directory. No images are provided out of the box, but it's possible to add your own by saving your images using the file name scheme found at the top of [game.go](game.go):

```
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
```