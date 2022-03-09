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
		Title:     "Balls n' Dongs",
		Bounds:    pixel.R(0, 0, 1024, 1024),
		VSync:     true,
		Resizable: true,
	}
	var App *app
	App = &app{sMenu, pixelg.New(cfg), buildMenuHandler(func() { App.state = sGame }), buildGameManager()}
	last := time.Now()
	for !App.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		App.HandleInput(dt)
		App.Draw()
		//<-fps
	}
}

func main() {
	pixelgl.Run(run)
}
