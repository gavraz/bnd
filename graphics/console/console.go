package console

import (
	"fmt"
	term "github.com/nsf/termbox-go"
)

type choicer interface {
	CurrentChoice() int
	Choices() []string
}

type menuHandler interface {
	NextChoice()
	PrevChoice()
	Choose()
	GoBack()
}

type Handler struct {
	numOfLines int
}

func New(numOfLines int) *Handler {
	return &Handler{numOfLines: numOfLines}
}

func (h *Handler) clearScreen() {
	for i := 0; i < h.numOfLines; i++ {
		fmt.Println()
	}
}

func (h *Handler) DrawMenu(c choicer) {
	h.clearScreen()

	current := c.CurrentChoice()
	for i, item := range c.Choices() {
		if i == current {
			fmt.Print(">")
		} else {
			fmt.Print(" ")
		}
		fmt.Println(item)
	}
	fmt.Println()
}

func (h *Handler) HandleInput(menuHandler menuHandler) {
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
