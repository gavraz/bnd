package main

import (
	"bnd/menu"
	"bnd/pixelDisplay"
	"fmt"
	"time"

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

func pixelHandleInput(menuHandler *menu.Handler, win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyS) || win.JustPressed(pixelgl.KeyDown) {
		menuHandler.NextChoice()
	}
	if win.JustPressed(pixelgl.KeyW) || win.JustPressed(pixelgl.KeyUp) {
		menuHandler.PrevChoice()
	}
	if win.JustPressed(pixelgl.KeyEnter) {
		menuHandler.Choose()
	}
	if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyBackspace) {
		menuHandler.GoBack()
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Balls n' Dongs",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	menuHandler := buildMenuHandler()
	displayHandler := pixelDisplay.New()

	fps := time.Tick(time.Second / 60)

	for !win.Closed() {
		displayHandler.DrawMenu(menuHandler, win)
		pixelHandleInput(menuHandler, win)

		win.Update()

		<-fps
	}
}

func main() {
	fmt.Println("Client: Hello B&D")

	err := term.Init()
	if err != nil {
		panic(err)
	}

	pixelgl.Run(run)

	defer term.Close()

	/*
		menuHandler := buildMenuHandler()
		displayHandler := display.New(5)

		for {
			displayHandler.DrawMenu(menuHandler)
			handleInput(menuHandler)
		}
	*/

}
