package screen

import (
	"time"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// loadingScreen is a screen that extends the Screen interface
type loadingScreen struct {
	Screen
}

func (ls *loadingScreen) Update() {
	rl.ClearBackground(rl.RayWhite)
	gui.Label(rl.NewRectangle(10, 10, 310, 25), "Loading...")
	time.Sleep(10 * time.Second)
	Manager.SetScreen(&menuScreen{})
}

func (ls *loadingScreen) OnEnter() {
	// Enter the loading screen
}

func (ls *loadingScreen) OnExit() {
	// Exit the loading screen
}
