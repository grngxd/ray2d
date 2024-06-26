package main

import (
	"fmt"
	"ray2d/networking"
	"ray2d/screen"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var c *networking.Connection

func main() {
	fmt.Println("Hello, World!")
	rl.InitWindow(800, 450, "RAY2D")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		// rl.ClearBackground(rl.RayWhite)
		// gui.Label(rl.NewRectangle(10, 10, 310, 25), "connect to server:")
		// serverAddress := "localhost:8080"
		// gui.TextBox(rl.NewRectangle(10, 40, 310, 25), &serverAddress, 8, true)
		// if gui.Button(rl.NewRectangle(10, 70, 310, 25), "Connect") {
		// 	c = networking.NewConnection(serverAddress)
		// 	if c == nil {
		// 		fmt.Println("Failed to connect to server")
		// 		continue
		// 	}

		// 	go c.Connect()
		// }

		// if gui.Button(rl.NewRectangle(10, 100, 310, 25), "Disconnect") {
		// 	if c != nil {
		// 		go c.Disconnect()
		// 	}
		// }

		screen.Manager.Update()

		rl.EndDrawing()
	}
}
