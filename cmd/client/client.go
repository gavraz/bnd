package main

import (
	"bnd/menu"
	"fmt"
	term "github.com/nsf/termbox-go"
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

func main() {
	fmt.Println("Client: Hello B&D")

	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

	menuHandler := buildMenuHandler()

	for {
		handleInput(menuHandler)
		drawMenu(menuHandler)
	}

}
