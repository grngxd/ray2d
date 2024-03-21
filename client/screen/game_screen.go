package screen

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// loadingScreen is a screen that extends the Screen interface
type gameScreen struct {
	Screen
}

func (ls *gameScreen) Update() {
	rl.ClearBackground(rl.DarkBrown)
	gui.Label(rl.NewRectangle(10, 10, 310, 25), "GAME!!!")
}

func (ls *gameScreen) OnEnter() {
	// Enter the loading screen
}

func (ls *gameScreen) OnExit() {
	// Exit the loading screen
}
