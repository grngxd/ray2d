package screen

import (
	"fmt"
	"ray2d/networking"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// loadingScreen is a screen that extends the Screen interface
type menuScreen struct {
	Screen
}

var c *networking.Connection

func (ls *menuScreen) Update() {
	rl.ClearBackground(rl.RayWhite)
	gui.Label(rl.NewRectangle(10, 10, 310, 25), "connect to server:")
	serverAddress := "localhost:8080"
	gui.TextBox(rl.NewRectangle(10, 40, 310, 25), &serverAddress, 8, true)
	if gui.Button(rl.NewRectangle(10, 70, 310, 25), "Connect") {
		c = networking.NewConnection(serverAddress)
		if c == nil {
			fmt.Println("Failed to connect to server")
			return
		}

		go c.Connect()
	}

	if gui.Button(rl.NewRectangle(10, 100, 310, 25), "Disconnect") {
		if c != nil {
			go c.Disconnect()
		}
	}
}

func (ls *menuScreen) OnEnter() {
	// Enter the loading screen
}

func (ls *menuScreen) OnExit() {
	// Exit the loading screen
}
