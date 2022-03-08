package main

import (
	pixelg "bnd/graphics/pixel"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	fmt.Println("Client: Hello B&D")

	cfg := pixelgl.WindowConfig{
		Title:     "Balls n' Dongs",
		Bounds:    pixel.R(0, 0, 1024, 1024),
		VSync:     true,
		Resizable: true,
	}

	displayHandler := pixelg.New()
	displayHandler.Init(cfg)

	gameManager := buildGameManager()
	//menuHandler := buildMenuHandler()
	//fps := time.Tick(time.Second / 30) // Test out if physics are working as intended with various fps values
	last := time.Now()
	for !displayHandler.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		//displayHandler.DrawMenu(menuHandler)
		//displayHandler.HandleInput(menuHandler)
		displayHandler.DrawGame(gameManager)
		displayHandler.HandleInput(gameManager)
		gameManager.Update(dt)
		displayHandler.Update()
		//<-fps
	}
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

/*
.\client.exe -cpuprofile=cpuprof.prof
go tool pprof -http=":8000" pprofbin .\cpuprof
*/

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Fatal("could not close file ", err)
			}
		}(f) // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// Your program here
	pixelgl.Run(run)
}
