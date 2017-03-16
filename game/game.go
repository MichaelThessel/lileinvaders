package game

import (
	"github.com/MichaelThessel/spacee/app"
	"github.com/veandco/go-sdl2/sdl"
)

// Game holds the game state
type Game struct {
	a  *app.App
	p  *player
	bl *bulletList
}

// New returns a new game
func New(a *app.App) *Game {
	g := &Game{a: a}

	g.p = newPlayer(a.GetRenderer())
	g.bl = &bulletList{}

	g.setup()

	return g
}

// setup sets up the game
func (g *Game) setup() {
	g.a.RegisterKeyCallback(sdl.K_LEFT, func() { g.p.Move('l') })
	g.a.RegisterKeyCallback(sdl.K_RIGHT, func() { g.p.Move('r') })

	g.a.RegisterKeyCallback(sdl.K_SPACE, func() { g.p.Fire(g.bl) })

	g.a.RegisterRenderCallback(1, g.p.Draw)
	g.a.RegisterRenderCallback(1, g.bl.Draw)
}
