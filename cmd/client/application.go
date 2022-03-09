package main

import (
	"bnd/game"
	pixelg "bnd/graphics/pixel"
	"bnd/input"
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
	state            state
	displayHandler   *pixelg.Handler
	mainMenuHandler  *menu.Handler
	pauseMenuHandler *menu.Handler
	gameManager      *game.Manager
	inputController  *input.Controller
}

func (a *application) Init() {
	cfg := pixelgl.WindowConfig{
		Title:     "Balls n' Dongs",
		Bounds:    pixel.R(0, 0, 1920, 1080),
		VSync:     true,
		Resizable: true,
	}

	a.state = stateMenu
	a.displayHandler = pixelg.NewHandler()
	a.displayHandler.Init(cfg)
	a.mainMenuHandler = buildMainMenuHandler(a.StartGame, a.changeResolution)
	a.pauseMenuHandler = buildPauseMenuHandler(
		a.ResumeGame,
		a.BackToMenu,
		a.QuitGame,
		a.RestartGame,
		a.changeResolution)
	a.gameManager = buildGameManager()
	a.inputController = input.NewController()
}

func (a *application) HandleInput() {
	if a.state == stateMenu {
		a.inputController.HandleMenuInput(a.displayHandler.JustPressed, a.mainMenuHandler)
	} else if a.state == statePause {
		a.inputController.HandleMenuInput(a.displayHandler.JustPressed, a.pauseMenuHandler)
	} else if a.state == stateGame {
		a.inputController.HandleGameInput(a.displayHandler.Pressed, a.PauseGame, a.gameManager.MovePlayer)
	}
}

func (a *application) Update(dt float64) {
	if a.state == stateGame {
		a.gameManager.Update(dt)
	}
	a.displayHandler.Update()
}

func (a *application) Draw() {
	if a.state == stateMenu {
		a.displayHandler.DrawMenu(a.mainMenuHandler)
	} else if a.state == statePause {
		a.displayHandler.DrawMenu(a.pauseMenuHandler)
	} else if a.state == stateGame {
		a.displayHandler.DrawGame(a.gameManager)
	}

}

func (a *application) Running() bool {
	return a.displayHandler.Closed()
}

func (a *application) changeResolution(width, height int) {
	a.displayHandler.ChangeResolution(width, height)
}

func (a *application) BackToMenu() {
	a.state = stateMenu
}

func (a *application) PauseGame() {
	a.state = statePause
}

func (a *application) ResumeGame() {
	a.state = stateGame
}

func (a *application) RestartGame() {
	a.gameManager.ResetGame()
	a.state = stateGame
}

func (a *application) QuitGame() {
	os.Exit(0)
}

func (a *application) StartGame() {
	a.gameManager.ResetGame()
	a.state = stateGame
}
