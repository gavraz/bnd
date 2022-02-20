package main

import (
	"bnd/menu"
	"fmt"
	"os"
)

func buildMenuHandler() *menu.Handler {
	settings := (&menu.Builder{}).
		WithOption("Name", func() {
		}).
		WithOption("Character", func() {
		}).Build()

	mainMenu := (&menu.Builder{}).
		WithOption("Start", func() {
			fmt.Println("Starting game...")
		}).
		WithSubMenu("Settings", settings).
		WithOption("Quit", func() {
			fmt.Println("quitting")
			os.Exit(0)
		}).Build()

	return menu.NewHandler(mainMenu)
}
