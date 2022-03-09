package main

import (
	"bnd/game"
	pixelg "bnd/graphics/pixel"
	"github.com/gavraz/menu/menu"
)

type State int

const (
	sMenu State = iota
	sGame
)

type app struct {
	state          State
	displayHandler *pixelg.Handler
	menuHandler    *menu.Handler
	gameManager    *game.Manager
}

func (a *app) HandleInput(dt float64) {
	if a.state == sMenu {
		a.displayHandler.HandleMenuInput(a.menuHandler)
	} else {
		a.displayHandler.HandleInput(a.gameManager, func() { a.state = sMenu })
		a.gameManager.Update(dt)
	}
	a.displayHandler.Update()
}

func (a *app) Draw() {
	if a.state == sMenu {
		a.displayHandler.DrawMenu(a.menuHandler)
	} else {
		a.displayHandler.DrawGame(a.gameManager)
	}

}

func (a *app) Closed() bool {
	return a.displayHandler.Closed()
}
