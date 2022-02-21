package main

import (
	pixelg "bnd/graphics/pixel"
	"fmt"
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

	menuHandler := buildMenuHandler()

	for !displayHandler.Closed() {
		displayHandler.DrawMenu(menuHandler)
		displayHandler.HandleInput(menuHandler)

		displayHandler.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
