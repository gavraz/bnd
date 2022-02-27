package main

import (
	pixelg "bnd/graphics/pixel"
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	fmt.Println("Client: Hello B&D")

	cfg := pixelgl.WindowConfig{
		Title:  "Balls n' Dongs",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	displayHandler := pixelg.New(cfg)
	displayHandler.Init()

	gameManager := buildGameManager()
	//menuHandler := buildMenuHandler()
	//fps := time.Tick(time.Second / 30) // Test out if physics are working as intended with various fps values
	last := time.Now()
	for !displayHandler.Closed() {
		dt := time.Since(last).Seconds()
		//displayHandler.DrawMenu(menuHandler)
		//displayHandler.HandleInput(menuHandler)
		displayHandler.DrawGame(gameManager)
		last = time.Now()
		displayHandler.HandleInput(gameManager, dt)
		gameManager.Update(dt)
		displayHandler.Update()
		//<-fps
	}
}

func main() {
	pixelgl.Run(run)
}
