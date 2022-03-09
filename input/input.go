package input

import (
	"bnd/game"
)

type Key int

const (
	KeyW Key = iota
	KeyD
	KeyS
	KeyA
	KeyEsc
	KeyEnter
)

type action int

const (
	Up action = iota
	Down
	Left
	Right
	ESC
	Enter
)

type Controller struct {
	mappings map[action]Key
}

func NewController() *Controller {
	c := &Controller{
		mappings: map[action]Key{},
	}

	c.setDefault()
	return c
}

func (c *Controller) setDefault() {
	c.mappings[Up] = KeyW
	c.mappings[Down] = KeyS
	c.mappings[Left] = KeyA
	c.mappings[Right] = KeyD
	c.mappings[ESC] = KeyEsc
	c.mappings[Enter] = KeyEnter
}

type movePlayerFunc func(direction game.Direction)

func (c *Controller) HandleGameInput(isPressed func(key Key) bool, pauseGame func(), movePlayer movePlayerFunc) {
	var dir game.Direction

	if isPressed(c.mappings[Up]) {
		dir.Up()
	}
	if isPressed(c.mappings[Down]) {
		dir.Down()
	}
	if isPressed(c.mappings[Left]) {
		dir.Left()
	}
	if isPressed(c.mappings[Right]) {
		dir.Right()
	}
	if isPressed(c.mappings[ESC]) {
		pauseGame()
	}

	movePlayer(dir)
}

type menuHandler interface {
	NextChoice()
	PrevChoice()
	Choose()
	GoBack()
}

func (c *Controller) HandleMenuInput(isPressed func(key Key) bool, menuHandler menuHandler) {
	if isPressed(c.mappings[Down]) {
		menuHandler.NextChoice()
	}
	if isPressed(c.mappings[Up]) {
		menuHandler.PrevChoice()
	}
	if isPressed(c.mappings[Enter]) {
		menuHandler.Choose()
	}
	if isPressed(c.mappings[ESC]) {
		menuHandler.GoBack()
	}
}
