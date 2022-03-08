package main

import (
	"bnd/game"
	pixelg "bnd/graphics/pixel"
	"fmt"
	"github.com/gavraz/menu/menu"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type State int

const (
	Menu State = iota
	Game
)

type appManager struct {
	state          State
	displayHandler *pixelg.Handler
	menuHandler    *menu.Handler
	gameManager    *game.Manager
}

func New(dh *pixelg.Handler, menu *menu.Handler, gm *game.Manager) *appManager {
	return &appManager{Menu, dh, menu, gm}
}

// GetInput - Maybe
func (a *appManager) GetInput() []bool {
	return []bool{}
}

func (a *appManager) HandleInput(dt float64) {
	if a.state == Menu {
		a.displayHandler.HandleMenuInput(a.menuHandler)
	}
	if a.state == Game {
		a.displayHandler.HandleInput(a.gameManager, dt)
		a.gameManager.Update(dt)
	}
	a.displayHandler.Update()
}

func (a *appManager) Draw() {
	if a.state == Menu {
		a.displayHandler.DrawMenu(a.menuHandler)
	}
	if a.state == Game {
		a.displayHandler.DrawGame(a.gameManager)
	}

}

func run() {
	fmt.Println("Client: Hello B&D")

	cfg := pixelgl.WindowConfig{
		Title:     "Balls n' Dongs",
		Bounds:    pixel.R(0, 0, 1024, 1024),
		VSync:     false,
		Resizable: true,
	}
	var App *appManager
	App = New(pixelg.New(), buildMenuHandler(func() { App.state = Game }), buildGameManager())
	App.displayHandler.Init(cfg)
	last := time.Now()
	var count, total int
	for !App.displayHandler.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		App.HandleInput(dt)
		App.Draw()
		//<-fps
		fps := int(1.0 / dt)
		if fps < 60 {
			count++
		}
		total++
	}
	fmt.Println("<60fps:", count, float64(count)/float64(total))
}

func main() {
	pixelgl.Run(run)
}
