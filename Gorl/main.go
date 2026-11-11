package main

import (
	"gorl/game"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(1280, 720, "gorl")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.SetWindowState(rl.FlagWindowResizable)

	game.Load()
	defer game.Unload()

	frameTime := 0
	for !rl.WindowShouldClose() {
		frameStart := time.Now()
		
		rl.BeginDrawing()

		game.Loop(frameTime)

		rl.EndDrawing()

		frameTime = int(time.Since(frameStart))
	}
}
