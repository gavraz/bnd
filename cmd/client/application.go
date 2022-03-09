package main

import (
	"bnd/game"
	pixelg "bnd/graphics/pixel"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gavraz/menu/menu"
	"os"
)

type state int

const (
	stateMenu state = iota
	statePause
	stateGame
)

type application struct {
	appState         state
	displayHandler   *pixelg.Handler
	mainMenuHandler  *menu.Handler
	pauseMenuHandler *menu.Handler
	gameManager      *game.Manager
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
	a.mainMenuHandler = buildMainMenuHandler(func() { a.appState = stateGame }, a.changeResolution)
	a.pauseMenuHandler = buildPauseMenuHandler(func() { a.appState = stateGame },
		func() { a.appState = stateMenu },
		func() { os.Exit(0) },
		a.changeResolution)
	a.gameManager = buildGameManager()
}

func (a *application) HandleInput() {
	if a.appState == stateMenu {
		a.displayHandler.HandleMenuInput(a.mainMenuHandler)
	} else if a.appState == statePause {
		a.displayHandler.HandleMenuInput(a.pauseMenuHandler)
	} else if a.appState == stateGame {
		a.displayHandler.HandleGameInput(a.gameManager, a)
	}
}

func (a *application) Update(dt float64) {
	a.gameManager.Update(dt)
	a.displayHandler.Update()
}

func (a *application) Draw() {
	if a.appState == stateMenu {
		a.displayHandler.DrawMenu(a.mainMenuHandler)
	} else if a.appState == statePause {
		a.displayHandler.DrawMenu(a.pauseMenuHandler)
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
	a.gameManager.ResetGame()
	a.appState = stateMenu
}

func (a *application) PauseGame() {
	a.appState = statePause
}
