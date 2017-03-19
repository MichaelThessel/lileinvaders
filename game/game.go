package game

import (
	"github.com/MichaelThessel/spacee/app"
	"github.com/veandco/go-sdl2/sdl"
)

// Game holds the game state
type Game struct {
	c     *Config
	a     *app.App
	p     *player     // Player
	pbl   *bulletList // Player bullet list
	abl   *bulletList // Alien bullet list
	ag    *alienGrid  // Alien grid
	stats *stats      // Game stats
	score int         // Game score
}

// Config holds game configuration
type Config struct {
	agc *alienGridConfig
	pc  *playerConfig
}

// New returns a new game
func New(a *app.App) (*Game, error) {
	var err error

	g := &Game{
		a:     a,
		score: 0,
		pbl:   &bulletList{},
		abl:   &bulletList{},
		ag:    &alienGrid{},
	}
	g.initConfig()

	// Player
	g.p, err = newPlayer(a.GetRenderer(), g.c.pc)
	if err != nil {
		return nil, err
	}

	// Start a new level
	err = g.startLevel(a.GetRenderer())
	if err != nil {
		return nil, err
	}

	// Stats
	g.stats, err = newStats(a.GetRenderer(), g.c.pc.lifes)
	if err != nil {
		return nil, err
	}

	g.setup()

	return g, nil
}

// initConfig initalizes gthe game config
func (g *Game) initConfig() {
	g.c = &Config{
		agc: &alienGridConfig{
			rows:        5,
			cols:        10,
			marginRow:   20,
			marginCol:   20,
			returnPoint: 30,
			speed:       4,
			speedStep:   5,
			bulletSpeed: 15,
			fireRate:    0.05,
		},
		pc: &playerConfig{
			stepSize:    30,
			bulletSpeed: 30,
			lifes:       5,
		},
	}
}

// setup sets up the game
func (g *Game) setup() {
	// Keyboard
	g.a.RegisterKeyCallback(sdl.K_LEFT, func() { g.p.Move('l') })    // left
	g.a.RegisterKeyCallback(sdl.K_RIGHT, func() { g.p.Move('r') })   // right
	g.a.RegisterKeyCallback(sdl.K_SPACE, func() { g.p.Fire(g.pbl) }) // fire

	// Draw player
	g.a.RegisterRenderCallback(1, g.p.Draw)

	// Draw player & alien bullets
	g.a.RegisterRenderCallback(1, g.abl.Draw)
	g.a.RegisterRenderCallback(1, g.pbl.Draw)

	// Draw alien grid
	g.a.RegisterRenderCallback(1, g.ag.Draw)

	// Draw stats
	g.a.RegisterRenderCallback(1, func() { g.stats.Draw(g.p.lifes, g.score) })

	// Test if bullets have hit
	g.a.RegisterRenderCallback(1, func() {
		if hit, len := g.ag.testHit(g.pbl); hit {
			g.score += 30
			if len == 0 {
				g.startLevel(g.a.GetRenderer())
			}
		}
	})
	g.a.RegisterRenderCallback(1, func() { g.p.testHit(g.abl) })

	// Aliens fire
	g.a.RegisterRenderCallback(1, func() { g.ag.fire(g.abl) })
}

// startLevel starts a level
func (g *Game) startLevel(r *sdl.Renderer) error {
	// Reset alien grid
	ag, err := newAlienGrid(r, g.c.agc)
	*g.ag = *ag

	// Reset bullet list
	bl := bulletList{}
	*g.abl = bl
	*g.pbl = bl

	return err
}
