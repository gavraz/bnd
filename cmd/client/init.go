package main

import (
	"bnd/game"
	"fmt"
	"os"

	"github.com/gavraz/menu/menu"
)

func buildMenuHandler(startGame func(), ChangeResolution func(int, int)) *menu.Handler {
	h := menu.NewHandler()

	resolution := menu.NewBuilder(h).
		WithOption("1920x1080", func() {
			ChangeResolution(1920, 1080)
		}).
		WithOption("1280x720", func() {
			ChangeResolution(1280, 720)
		}).
		WithOption("800x600", func() {
			ChangeResolution(800, 600)
		}).
		WithOption("640x480", func() {
			ChangeResolution(640, 480)
		}).
		WithOption("320x240", func() {
			ChangeResolution(320, 240)
		}).
		WithGoBack("Go Back").
		Build()

	graphics := menu.NewBuilder(h).
		WithSubMenu("Change Resolution", resolution).
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
