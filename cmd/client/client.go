package main

import (
	"bnd/menu"
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	term "github.com/nsf/termbox-go"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func drawMenu(h *menu.Handler) {
	current := h.CurrentChoice()
	for i, item := range h.Choices() {
		if i == current {
			fmt.Print(">")
		} else {
			fmt.Print(" ")
		}
		fmt.Println(item)
	}
	fmt.Println()
	fmt.Println()
}

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
	cfg := pixelgl.WindowConfig{
		Title:  "Balls n' Dongs",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(100, 500), basicAtlas)

	win.Clear(colornames.Skyblue)
	fmt.Fprintln(basicTxt, "Hello, text!")
	fmt.Fprintln(basicTxt, "I support multiple lines!")
	fmt.Fprintf(basicTxt, "And I'm an %s, yay!", "io.Writer")

	for !win.Closed() {
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 4))
		win.Update()
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

	menuHandler := buildMenuHandler()

	for {
		handleInput(menuHandler)
		drawMenu(menuHandler)
	}

}
