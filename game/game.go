package game

import (
	"fmt"

	"github.com/MichaelThessel/spacee/app"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	// Game scene constants
	sceneStart = "start"
	scenePlay  = "play"
	sceneEnd   = "end"
)

// Game holds the game state
type Game struct {
	c     *Config
	a     *app.App
	scene string
	start *start      // Start screen
	end   *end        // End screen
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
	g := &Game{
		scene: sceneStart,
		a:     a,
		score: 0,
		pbl:   &bulletList{},
		abl:   &bulletList{},
		ag:    &alienGrid{},
	}
	g.initConfig()

	if err := g.switchScene(sceneStart); err != nil {
		return nil, err
	}

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
			speedMax:    5,
			speedStep:   3,
			bulletSpeed: 15,
			fireRate:    0.05,
			stepSizeX:   10,
			stepSizeY:   10,
		},
		pc: &playerConfig{
			stepSize:    30,
			bulletSpeed: 30,
			lifes:       5,
		},
	}
}

// switchScene switches to a different scene
func (g *Game) switchScene(scene string) error {
	g.a.ClearCallbacks()

	switch scene {
	case sceneStart:
		g.scene = sceneStart
		return g.sceneStart()
	case scenePlay:
		g.scene = scenePlay
		return g.scenePlay()
	case sceneEnd:
		g.scene = sceneEnd
		return g.sceneEnd()
	default:
		panic(fmt.Sprintf("Invalid scene %s", scene))
	}
}

// sceneStart sets up the start screen
func (g *Game) sceneStart() error {
	// Start screen
	var err error
	g.start, err = newStart(g.a.GetRenderer())
	if err != nil {
		return err
	}

	// Draw start screen
	g.a.RegisterRenderCallback(1, g.start.Draw)

	g.a.RegisterKeyCallback(sdl.K_RETURN, func() { g.switchScene(scenePlay) }) // start

	return nil
}

// scenePlay sets up the game
func (g *Game) scenePlay() error {
	// Player
	var err error
	g.p, err = newPlayer(g.a.GetRenderer(), g.c.pc)
	if err != nil {
		return err
	}

	// Start a new level
	err = g.startLevel(g.a.GetRenderer())
	if err != nil {
		return err
	}

	// Stats
	g.score = 0
	g.stats, err = newStats(g.a.GetRenderer(), g.c.pc.lifes)
	if err != nil {
		return err
	}

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

	// Test if player bullets have hit
	g.a.RegisterRenderCallback(1, func() {
		if hit, len := g.ag.testHit(g.pbl); hit {
			g.score += 30
			if len == 0 {
				g.startLevel(g.a.GetRenderer())
			}
		}
	})

	// Test if alien bullets have hit
	g.a.RegisterRenderCallback(1, func() {
		if g.p.testHit(g.abl) {
			g.switchScene(sceneEnd)
		}
	})

	// Test if aliens have reached the ground
	g.a.RegisterRenderCallback(1, func() {
		if g.ag.testBoundary() {
			g.switchScene(sceneEnd)
		}
	})

	// Test if aliens collided with player
	g.a.RegisterRenderCallback(1, func() {
		if g.ag.testPlayerCollission(g.p) {
			g.switchScene(sceneEnd)
		}
	})

	// Aliens fire
	g.a.RegisterRenderCallback(1, func() { g.ag.fire(g.abl) })

	return nil
}

// sceneEnd sets up the end scene
func (g *Game) sceneEnd() error {
	// End screen
	var err error
	g.end, err = newEnd(g.a.GetRenderer(), g.score)
	if err != nil {
		return err
	}

	// Draw end screen
	g.a.RegisterRenderCallback(1, g.end.Draw)

	g.a.RegisterKeyCallback(sdl.K_RETURN, func() { g.switchScene(scenePlay) }) // start

	return nil
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
