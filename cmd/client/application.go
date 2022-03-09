package main

import (
	"bnd/game"
	pixelg "bnd/graphics/pixel"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
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

func NewApplication() *application {
	return &application{}
}

func (a *application) Init() {
	cfg := pixelgl.WindowConfig{
		Title:     "Balls n' Dongs",
		Bounds:    pixel.R(0, 0, 1920, 1080),
		VSync:     true,
		Resizable: true,
	}

	a.appState = stateMenu
	a.displayHandler = pixelg.NewHandler()
	a.displayHandler.Init(cfg)
	a.menuHandler = buildMenuHandler(func() { a.appState = stateGame }, a.changeResolution)
	a.gameManager = buildGameManager()
}

func (a *application) HandleInput() {
	if a.appState == stateMenu {
		a.displayHandler.HandleMenuInput(a.menuHandler)
	} else if a.appState == stateGame {
		a.displayHandler.HandleGameInput(a.gameManager, a, a.gameManager)
	}
}

func (a *application) Update(dt float64) {
	a.gameManager.Update(dt)
	a.displayHandler.Update()
}

func (a *application) Draw() {
	if a.appState == stateMenu {
		a.displayHandler.DrawMenu(a.menuHandler)
	} else if a.appState == stateGame {
		a.displayHandler.DrawGame(a.gameManager)
	}

}

func (a *application) Closed() bool {
	return a.displayHandler.Closed()
}

func (a *application) changeResolution(width, height int) {
	a.displayHandler.ChangeResolution(width, height)
}

func (a *application) BackToMenu() {
	a.appState = stateMenu
}
