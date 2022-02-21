package main

import (
	pixelg "bnd/graphics/pixel"
	"bnd/menu"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	term "github.com/nsf/termbox-go"
)

func handleInput(menuHandler *menu.Handler) {
	switch ev := term.PollEvent(); ev.Key {
	case term.KeyArrowDown:
		menuHandler.NextChoice()
	case term.KeyArrowUp:
		menuHandler.PrevChoice()
	case term.KeyEnter:
		menuHandler.Choose()
	case term.KeySpace:
		menuHandler.GoBack()
	}
}

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
