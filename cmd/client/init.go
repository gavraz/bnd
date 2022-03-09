package main

import (
	"bnd/game"
	"fmt"
	"os"

	"github.com/gavraz/menu/menu"
)

func buildMainMenuHandler(startGame func(), changeResolution func(int, int)) *menu.Handler {
	h := menu.NewHandler()

	settings := buildSettings(h, changeResolution)

	mainMenu := menu.NewBuilder(h).
		WithOption("Start", func() {
			startGame()
		}).
		WithSubMenu("Settings", settings).
		WithOption("Quit", func() {
			fmt.Println("quitting")
			os.Exit(0)
		}).
		Build()

	h.SetMenu(mainMenu)

	return h
}

func buildPauseMenuHandler(resumeGame func(), toMainMenu func(), quitGame func(), restartGame func(), changeResolution func(int, int)) *menu.Handler {
	h := menu.NewHandler()

	pauseMenu := menu.NewBuilder(h).
		WithOption("Resume", func() {
			resumeGame()
		}).
		WithOption("Restart Game", func() {
			restartGame()
		}).
		WithSubMenu("Settings", buildSettings(h, changeResolution)).
		WithOption("Main Menu", func() {
			toMainMenu()
		}).
		WithOption("Quit", func() {
			quitGame()
		}).
		Build()

	h.SetMenu(pauseMenu)

	return h
}

func buildSettings(h *menu.Handler, changeResolution func(int, int)) *menu.Menu {
	resolution := menu.NewBuilder(h).
		WithOption("1920x1080", func() {
			changeResolution(1920, 1080)
		}).
		WithOption("1280x720", func() {
			changeResolution(1280, 720)
		}).
		WithOption("800x600", func() {
			changeResolution(800, 600)
		}).
		WithOption("640x480", func() {
			changeResolution(640, 480)
		}).
		WithGoBack("Go Back").
		Build()

	graphics := menu.NewBuilder(h).
		WithSubMenu("Resolution", resolution).
		WithGoBack("Go Back").
		Build()

	settings := menu.NewBuilder(h).
		WithSubMenu("Graphics", graphics).
		WithOption("Name", func() {
		}).
		WithOption("Character", func() {
		}).
		WithGoBack("Go Back").
		Build()
	return settings
}

func buildGameManager() *game.Manager {
	m := game.NewManager()

	return m
}
