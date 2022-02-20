package main

import (
	"bnd/display"
	"bnd/menu"
	"fmt"
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

func main() {
	fmt.Println("Client: Hello B&D")

	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

	menuHandler := buildMenuHandler()
	displayHandler := display.New(5)

	for {
		displayHandler.DrawMenu(menuHandler)
		handleInput(menuHandler)
	}

}
