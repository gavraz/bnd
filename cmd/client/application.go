package main

import (
	"bnd/game"
	pixelg "bnd/graphics/pixel"
	"github.com/gavraz/menu/menu"
)

type state int

const (
	stateMenu state = iota
	stateGame
)

type application struct {
	appState       state
	displayHandler *pixelg.Handler
	menuHandler    *menu.Handler
	gameManager    *game.Manager
}

func (a *application) HandleInput(dt float64) {
	if a.appState == stateMenu {
		a.displayHandler.HandleMenuInput(a.menuHandler)
	} else {
		a.displayHandler.HandleInput(a.gameManager, func() { a.appState = stateMenu })
		a.gameManager.Update(dt)
	}
	a.displayHandler.Update()
}

func (a *application) Draw() {
	if a.appState == stateMenu {
		a.displayHandler.DrawMenu(a.menuHandler)
	} else {
		a.displayHandler.DrawGame(a.gameManager)
	}

}

func (a *application) Closed() bool {
	return a.displayHandler.Closed()
}
