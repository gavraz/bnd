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
	m.Add("current-player", game.NewPlayer(&game.GObject{
		Center: game.Vector2{
			X: 0,
			Y: 0,
		},
		BaseSpeed:     3,
		CollisionType: game.Circle,
		Width:         0.05,
		Height:        0.05,
		Mass:          100,
	}, 100))
	m.Add("crate", &game.Crate{
		Object: &game.GObject{
			Center: game.Vector2{
				X: -0.2,
				Y: -0.2,
			},
			CollisionType: game.Rectangle,
			Width:         0.1,
			Height:        0.1,
			Mass:          1,
		},
	})
	m.Add("crate2", &game.Crate{
		Object: &game.GObject{
			Center: game.Vector2{
				X: -0.3,
				Y: -0.3,
			},
			CollisionType: game.Rectangle,
			Width:         0.1,
			Height:        0.1,
			Mass:          1,
		},
	})
	m.Add("bouncing-ball", &game.BouncingBall{
		Object: &game.GObject{
			Center: game.Vector2{
				X: 0.1,
				Y: 0.3,
			},
			CollisionType: game.Circle,
			Width:         0.1,
			Height:        0.1,
			Mass:          1,
		},
	})

	return m
}
