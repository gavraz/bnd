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
			X: 500,
			Y: 500,
		},
		BaseSpeed:     1024,
		CollisionType: game.Circle,
		Width:         30,
		Height:        30,
		Mass:          100,
	}, 100))
	m.Add("crate", &game.Crate{
		Object: &game.GObject{
			Center: game.Vector2{
				X: 200,
				Y: 200,
			},
			CollisionType: game.Rectangle,
			Width:         50,
			Height:        50,
			Mass:          1,
		},
	})
	m.Add("crate2", &game.Crate{
		Object: &game.GObject{
			Center: game.Vector2{
				X: 300,
				Y: 300,
			},
			CollisionType: game.Rectangle,
			Width:         50,
			Height:        50,
			Mass:          1,
		},
	})
	m.Add("bouncing-ball", &game.BouncingBall{
		Object: &game.GObject{
			Center: game.Vector2{
				X: 100,
				Y: 300,
			},
			CollisionType: game.Circle,
			Width:         50,
			Height:        50,
			Mass:          1,
		},
	})
	// m.Add("bottom-wall", &game.Wall{
	// 	Object: &game.GObject{
	// 		Center: game.Vector2{
	// 			X: 500,
	// 			Y: 600,
	// 		},
	// 		CollisionType:   game.Rectangle,
	// 		Width:           50,
	// 		Height:          50,
	// 		CollisionWidth:  50,
	// 		CollisionHeight: 50,
	// 	},
	// })

	return m
}
