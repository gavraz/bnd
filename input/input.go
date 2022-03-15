package input

import (
	"bnd/game"
	"time"
)

type Key int

const (
	KeyW Key = iota
	KeyD
	KeyS
	KeyA
	KeyEsc
	KeyEnter
	KeySpace
)

type action int

const (
	playerUp action = iota
	playerDown
	playerLeft
	playerRight
	playerAbility
	menuUp
	menuDown
	menuBack
	menuEnter
	pauseGame
)

type Controller struct {
	actionToKey map[action]Key
}

func NewController() *Controller {
	c := &Controller{
		actionToKey: map[action]Key{},
	}

	c.setDefault()
	return c
}

func (c *Controller) setDefault() {
	c.actionToKey[playerUp] = KeyW
	c.actionToKey[playerDown] = KeyS
	c.actionToKey[playerLeft] = KeyA
	c.actionToKey[playerRight] = KeyD
	c.actionToKey[menuUp] = KeyW
	c.actionToKey[menuDown] = KeyS
	c.actionToKey[menuBack] = KeyEsc
	c.actionToKey[menuEnter] = KeyEnter
	c.actionToKey[pauseGame] = KeyEsc
	c.actionToKey[playerAbility] = KeySpace
}

type movePlayerFunc func(direction game.MoveDirection)

func (c *Controller) HandleGameInput(isPressed func(key Key) bool, justPressed func(key Key) bool, pause func(), movePlayer movePlayerFunc, attack func(), dt float64) {
	var dir game.MoveDirection

	if isPressed(c.actionToKey[playerUp]) {
		dir.Up()
	}
	if isPressed(c.actionToKey[playerDown]) {
		dir.Down()
	}
	if isPressed(c.actionToKey[playerLeft]) {
		dir.Left()
	}
	if isPressed(c.actionToKey[playerRight]) {
		dir.Right()
	}
	if isPressed(c.actionToKey[pauseGame]) {
		pause()
	}
	if justPressed(c.actionToKey[playerAbility]) {
		attack()
		go func() {
			for i := 0; i < 1; i++ {

				time.Sleep(100 * time.Millisecond)
			}
		}()
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
	if isPressed(c.actionToKey[menuDown]) {
		menuHandler.NextChoice()
	}
	if isPressed(c.actionToKey[menuUp]) {
		menuHandler.PrevChoice()
	}
	if isPressed(c.actionToKey[menuEnter]) {
		menuHandler.Choose()
	}
	if isPressed(c.actionToKey[menuBack]) {
		menuHandler.GoBack()
	}
}
