package main

import (
	"fmt"
	"os"

	"github.com/MichaelThessel/spacee/app"
	"github.com/MichaelThessel/spacee/game"
)

func main() {
	// TODO: config needs to come from flags
	config := &app.Config{
		Width:     800,
		Height:    600,
		Title:     "e-Space",
		FrameRate: 30,
	}

	a, err := app.New(config)
	if err != nil {
		panic(fmt.Sprintf("couldn't set up window %v", err))
	}
	defer a.Destroy()

	g := game.New(a)
	g.Run()

	os.Exit(a.Run())
}
