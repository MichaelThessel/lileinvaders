package main

import (
	"fmt"
	"os"

	"github.com/MichaelThessel/spacee/app"
)

func main() {
	a := app.New()
	if err := a.Setup(); err != nil {
		panic(fmt.Sprintf("couldn't set up window %v", err))
	}
	defer a.Destroy()

	os.Exit(a.Run())
}
