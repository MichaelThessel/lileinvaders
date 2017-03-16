package game

import (
	"github.com/MichaelThessel/spacee/app"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	a *app.App
	p *player
}

func New(a *app.App) *Game {
	g := &Game{a: a}

	g.p = newPlayer(a.GetRenderer())

	return g
}

func (g *Game) Run() {
	g.a.RegisterKeyCallback(sdl.K_LEFT, func() { g.p.Move('l') })
	g.a.RegisterKeyCallback(sdl.K_RIGHT, func() { g.p.Move('r') })

	g.a.RegisterRenderCallback(1, g.p.Draw)
}
