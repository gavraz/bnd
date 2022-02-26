package main

import (
	"bnd/game"
	"fmt"
	"os"

	"github.com/gavraz/menu/menu"
)

func buildMenuHandler() *menu.Handler {
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
			fmt.Println("Starting game...")
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
	player1 := game.NewPlayer(&game.GObject{
		Point:        game.Point{X: 100, Y: 100},
		Velocity:     game.Vector2{X: 0, Y: 0},
		Acceleration: game.Vector2{X: 0, Y: 0},
		Direction:    game.Vector2{X: 0, Y: 0},
	}, 100)
	m.Add("current-player", player1)

	return m
}
