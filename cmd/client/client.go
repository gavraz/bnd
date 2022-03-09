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
		Bounds:    pixel.R(0, 0, 1920, 1080),
		VSync:     true,
		Resizable: true,
	}
	dh := pixelg.New()
	dh.Init(cfg)
	var app *application
	app = &application{stateMenu, dh, buildMenuHandler(func() { app.appState = stateGame }, displayHandler.ChangeResolution), buildGameManager()}
	last := time.Now()
	for !app.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		app.HandleInput(dt)
		app.Draw()
		//<-fps
	}
}

func main() {
	pixelgl.Run(run)
}
