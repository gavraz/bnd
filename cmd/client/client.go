package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

func run() {
	fmt.Println("Client: Hello B&D")

	app := NewApplication()
	app.Init()
	last := time.Now()
	for !app.Running() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		app.HandleInput()
		app.Draw()
		app.Update(dt)
	}
}

func main() {
	pixelgl.Run(run)
}
