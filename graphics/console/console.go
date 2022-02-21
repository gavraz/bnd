package console

import (
	"fmt"
)

type choicer interface {
	CurrentChoice() int
	Choices() []string
}

type Handler struct {
	numOfLines int
}

func New(numOfLines int) Handler {
	return Handler{numOfLines: numOfLines}
}

func (d Handler) clearScreen() {
	for i := 0; i < d.numOfLines; i++ {
		fmt.Println()
	}
}

func (d Handler) DrawMenu(c choicer) {
	d.clearScreen()

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
