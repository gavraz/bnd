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

type appManager struct {
	state          State
	displayHandler *pixelg.Handler
	menuHandler    *menu.Handler
	gameManager    *game.Manager
}

func (a *appManager) HandleInput(dt float64) {
	if a.state == sMenu {
		a.displayHandler.HandleMenuInput(a.menuHandler)
	} else {
		a.displayHandler.HandleInput(a.gameManager, dt, func() { a.state = sMenu })
		a.gameManager.Update(dt)
	}
	a.displayHandler.Update()
}

func (a *appManager) Draw() {
	if a.state == sMenu {
		a.displayHandler.DrawMenu(a.menuHandler)
	} else {
		a.displayHandler.DrawGame(a.gameManager)
	}

}

func (a *appManager) Closed() bool {
	return a.displayHandler.Closed()
}
