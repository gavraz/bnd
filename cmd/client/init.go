package main

import (
	"bnd/game"
	"fmt"
	"os"

	"github.com/gavraz/menu/menu"
)

func buildMenuHandler(startGame func()) *menu.Handler {
	h := menu.NewHandler()
	settings := menu.NewBuilder(h).
		WithOption("Name", func() {
		}).
		WithOption("Character", func() {
		}).
		WithGoBack("Go Back").
		Build()

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

func buildGameManager() *game.Manager {
	m := game.NewManager()
	m.InitGame()

	return m
}
