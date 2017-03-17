package game

import (
	"github.com/MichaelThessel/spacee/app"
	"github.com/veandco/go-sdl2/sdl"
)

// Game holds the game state
type Game struct {
	a   *app.App
	p   *player
	pbl *bulletList
	ag  *alienGrid
}

// New returns a new game
func New(a *app.App) (*Game, error) {
	g := &Game{a: a}

	var err error
	g.p, err = newPlayer(a.GetRenderer())
	if err != nil {
		return nil, err
	}

	g.ag, err = newAlienGrid(a.GetRenderer())
	if err != nil {
		return nil, err
	}

	// Player bullet list
	g.pbl = &bulletList{}

	g.setup()

	return g, nil
}

// setup sets up the game
func (g *Game) setup() {
	g.a.RegisterKeyCallback(sdl.K_LEFT, func() { g.p.Move('l') })
	g.a.RegisterKeyCallback(sdl.K_RIGHT, func() { g.p.Move('r') })

	g.a.RegisterKeyCallback(sdl.K_SPACE, func() { g.p.Fire(g.pbl) })

	g.a.RegisterRenderCallback(1, g.p.Draw)
	g.a.RegisterRenderCallback(1, g.pbl.Draw)
	g.a.RegisterRenderCallback(1, g.ag.Draw)
}
