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
	for !app.Closed() {
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
