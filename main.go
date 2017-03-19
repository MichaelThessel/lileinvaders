package main

import (
	"fmt"
	"os"

	"github.com/MichaelThessel/spacee/app"
	"github.com/MichaelThessel/spacee/game"
	"github.com/veandco/go-sdl2/sdl"
	mix "github.com/veandco/go-sdl2/sdl_mixer"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

func main() {
	// TODO: config needs to come from flags
	config := &app.Config{
		Width:     1200,
		Height:    800,
		Title:     "e-Space",
		FrameRate: 30,
	}

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Printf("could not initialize sdl: %v", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		fmt.Printf("could not initialize ttf: %v", err)
		os.Exit(1)
	}
	defer ttf.Quit()

	if err := mix.Init(0); err != nil {
		fmt.Printf("could not initialize mixer: %v", err)
		os.Exit(1)
	}
	defer mix.Quit()

	sndfmt := uint16(mix.DEFAULT_FORMAT)
	if err := mix.OpenAudio(44100, sndfmt, 2, 1024); err != nil {
		fmt.Printf("could not initialize audio: %v", err)
		os.Exit(1)
	}

	a, err := app.New(config)
	if err != nil {
		fmt.Printf("couldn't set up window %v", err)
	}
	defer a.Destroy()

	if _, err := game.New(a); err != nil {
		fmt.Printf("couldn't create game %v", err)
		os.Exit(1)
	}

	os.Exit(a.Run())
}
