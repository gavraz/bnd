package main

import (
	"fmt"
	"github.com/faiface/pixel/pixelgl"
	"time"
)

func run() {
	fmt.Println("Client: Hello B&D")
	app := &application{}
	app.Init()
	last := time.Now()
	for app.Running() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		app.HandleInput(dt)
		app.Draw()
		app.Update(dt)
	}
}

func main() {
	pixelgl.Run(run)
}
