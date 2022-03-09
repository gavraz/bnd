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
	playerUp action = iota
	playerDown
	playerLeft
	playerRight
	menuUp
	menuDown
	esc
	enter
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
	c.mappings[playerUp] = KeyW
	c.mappings[playerDown] = KeyS
	c.mappings[playerLeft] = KeyA
	c.mappings[playerRight] = KeyD
	c.mappings[menuUp] = KeyW
	c.mappings[menuDown] = KeyS
	c.mappings[esc] = KeyEsc
	c.mappings[enter] = KeyEnter
}

type movePlayerFunc func(direction game.Direction)

func (c *Controller) HandleGameInput(isPressed func(key Key) bool, pauseGame func(), movePlayer movePlayerFunc) {
	var dir game.Direction

	if isPressed(c.mappings[playerUp]) {
		dir.Up()
	}
	if isPressed(c.mappings[playerDown]) {
		dir.Down()
	}
	if isPressed(c.mappings[playerLeft]) {
		dir.Left()
	}
	if isPressed(c.mappings[playerRight]) {
		dir.Right()
	}
	if isPressed(c.mappings[esc]) {
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
	if isPressed(c.mappings[menuDown]) {
		menuHandler.NextChoice()
	}
	if isPressed(c.mappings[menuUp]) {
		menuHandler.PrevChoice()
	}
	if isPressed(c.mappings[enter]) {
		menuHandler.Choose()
	}
	if isPressed(c.mappings[esc]) {
		menuHandler.GoBack()
	}
}
