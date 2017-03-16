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
		fmt.Printf("couldn't set up window %v", err)
	}
	defer a.Destroy()

	if _, err := game.New(a); err != nil {
		fmt.Printf("couldn't create game %v", err)
		os.Exit(1)
	}

	os.Exit(a.Run())
}
