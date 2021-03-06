package game

import (
	"fmt"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
	mix "github.com/veandco/go-sdl2/sdl_mixer"
)

// alien holds the alien state
type alien struct {
	r *sdl.Renderer
	t *sdl.Texture
	x int32
	y int32
	w int32
	h int32
}

// newAlien generates a alien
func newAlien(r *sdl.Renderer, t *sdl.Texture, x, y int32) *alien {
	a := &alien{
		r: r,
		t: t,
		w: 80,
		h: 86,
		x: x,
		y: y,
	}

	return a
}

// Draw draws the alien
func (a *alien) Draw() {
	a.r.Copy(a.t, nil, &sdl.Rect{X: a.x, Y: a.y, W: a.w, H: a.h})
}

// alienGridConfig holds the alien grid config
type alienGridConfig struct {
	rows        int     // Number of rows
	cols        int     // Number of columns
	marginRow   int     // Space between rows
	marginCol   int     // Space between columns
	returnPoint int32   // When to switch the x direction
	speedMax    int     // Grid movement max speed
	speedStep   int     // After how many drops to increase the speed
	bulletSpeed int32   // Speed of a bullet
	fireRate    float64 // Rate at that the aliens fire
	stepSizeX   int32   // Horizontal step size
	stepSizeY   int32   // Vertical step size
}

// alienGrid holds the alien grid state
type alienGrid struct {
	c            *alienGridConfig
	r            *sdl.Renderer
	t            *sdl.Texture
	sounds       map[string]*mix.Chunk
	alienList    []*alien   // List of all aliens
	alienGridPos [][]*alien // List of all alien grid positions
	direction    int32      // direction of x movement (1: left, -1: right)
	dropCount    int        // How often the grid moved down in y
	speed        int        // Grid movement speed
	moveCounter  int        // Counts how many moves have been requested
}

// newAlienGrid creates a new alien grid
func newAlienGrid(r *sdl.Renderer, c *alienGridConfig) (*alienGrid, error) {
	maxX, _, _ := r.GetRendererOutputSize()

	ag := &alienGrid{
		c:         c,
		r:         r,
		direction: 1,
		dropCount: 0,
		speed:     1,
	}

	var err error
	ag.t, err = img.LoadTexture(ag.r, "assets/alien.png")
	if err != nil {
		return nil, fmt.Errorf("couldn't create alien texture: %v", err)
	}

	textureWidth := 80 // TODO: get this dynamically
	textureHeight := 86

	// Initialize the alien grid
	startX := (maxX - (textureWidth+ag.c.marginCol)*ag.c.cols - ag.c.marginCol) / 2
	startY := 50
	currentX := startX
	currentY := startY
	ag.alienGridPos = make([][]*alien, ag.c.rows)
	for row := 0; row < ag.c.rows; row++ {
		ag.alienGridPos[row] = make([]*alien, ag.c.cols)
		for col := 0; col < ag.c.cols; col++ {
			a := newAlien(r, ag.t, int32(currentX), int32(currentY))
			currentX += textureWidth + ag.c.marginCol
			ag.alienList = append(ag.alienList, a)
			ag.alienGridPos[row][col] = a
		}
		currentX = startX
		currentY += textureHeight + ag.c.marginRow
	}

	// Set sounds
	ag.sounds = make(map[string]*mix.Chunk, 0)
	ag.sounds["hit"], err = mix.LoadWAV("assets/sounds/alienhit.wav")
	if err != nil {
		return nil, fmt.Errorf("couldn't load sound: %v", err)
	}

	return ag, nil
}

// Draw renders the alien grid
func (ag *alienGrid) Draw() {
	for _, a := range ag.alienList {
		a.Draw()
	}
	ag.move()
}

// move moves the alien grid left and right and down
func (ag *alienGrid) move() {
	// Update the alien grid only every x frames
	ag.moveCounter++
	if ag.speed < ag.c.speedMax &&
		ag.moveCounter%(ag.c.speedMax-ag.speed) != 0 {
		return
	}
	ag.moveCounter = 0

	// Viewport && grid dimensions
	maxX, _, _ := ag.r.GetRendererOutputSize()
	x1, _, x2, _ := ag.getDimensions()

	// Check if the grid hits the boundary
	moveY := x2+ag.c.stepSizeX*ag.direction >= int32(maxX)-ag.c.returnPoint ||
		x1+ag.c.stepSizeX*ag.direction <= ag.c.returnPoint

	if moveY {
		// Increase the speed over time
		ag.dropCount++
		if ag.dropCount%ag.c.speedStep == 0 {
			ag.speed++
		}

		// Switch direction
		ag.direction *= -1
	}

	// Move all aliens
	for _, a := range ag.alienList {
		if moveY {
			a.y += ag.c.stepSizeY
		} else {
			a.x += ag.direction * ag.c.stepSizeX

		}
	}
}

// getDimentsions returns the current alien grid rectangle coordinates
func (ag *alienGrid) getDimensions() (x1, y1, x2, y2 int32) {
	x1, y1, x2, y2 = 0, 0, 0, 0

	for _, a := range ag.alienList {
		if a.x < x1 || x1 == 0 {
			x1 = a.x
		}
		if a.x+a.w > x2 {
			x2 = a.x + a.w
		}
		if a.y < y1 || y1 == 0 {
			y1 = a.y
		}
		if a.y+a.h > y2 {
			y2 = a.y + a.h
		}
	}

	return
}

// testHit checks if a bullet has hit an alien in the grid
func (ag *alienGrid) testHit(bl *bulletList) (bool, int) {
	x1, y1, x2, _ := ag.getDimensions()

	for _, b := range *bl {
		// Exit if bullet is beyond grid dimensions
		if b.x < x1 || b.x+b.w > x2 || b.y+b.h < y1 {
			continue
		}

		// Check if alien has been hit
		for _, a := range ag.alienList {
			if a.y+a.h < b.y || a.x > b.x+b.w || a.x+a.w < b.x {
				continue
			}

			// Hit detected: remove alien & bullet
			ag.remove(a)
			bl.remove(b)

			ag.sounds["hit"].Play(0, 0)

			return true, len(ag.alienList)
		}
	}

	return false, len(ag.alienList)
}

// testBoundary checks if the aliens have reached the ground
func (ag *alienGrid) testBoundary() bool {
	_, _, _, y := ag.getDimensions()
	_, maxY, _ := ag.r.GetRendererOutputSize()

	return y >= int32(maxY)
}

// testPlayerCollission checks if an alien has hit the player
func (ag *alienGrid) testPlayerCollission(p *player) (collission bool) {
	_, _, _, y := ag.getDimensions()

	// Test if grid is higher than player
	if p.y > y {
		return
	}

	bottomAliens := ag.bottomAliens()

	for c := 0; c < ag.c.cols; c++ {
		if bottomAliens[c] == nil {
			continue
		}

		a := bottomAliens[c]

		// Test if alien is higher than player
		if a.y+a.h < p.y {
			continue
		}

		if a.x >= p.x && a.x+a.w <= p.x+p.w {
			return true
		}
	}

	return
}

// remove removes an alien from the grid
func (ag *alienGrid) remove(a *alien) {
	// Remove alien from alien list
	tmpAl := []*alien{}
	for _, ta := range ag.alienList {
		if ta != a {
			tmpAl = append(tmpAl, ta)
		}
	}
	ag.alienList = tmpAl

	// Remove alien from alien grid position
	for r := range ag.alienGridPos {
		for c, ta := range ag.alienGridPos[r] {
			if ta == a {
				ag.alienGridPos[r][c] = nil
				return
			}
		}
	}
}

// fire randomly fires a bullets
func (ag *alienGrid) fire(bullets *bulletList) {
	if rand.Float64() > ag.c.fireRate {
		return
	}

	// Lowest row of aliens fires
	bottomAliens := ag.bottomAliens()
	for c := 0; c < ag.c.cols; c++ {
		if rand.Float64() > ag.c.fireRate {
			continue
		}

		if bottomAliens[c] != nil {
			a := bottomAliens[c]
			newBullet(
				ag.r,
				bullets,
				a.x+a.w/2,
				a.y+a.h,
				&bulletConfig{
					speed:     ag.c.bulletSpeed,
					direction: 1,
					colorR:    0xF6,
					colorG:    0x25,
					colorB:    0x9B,
				},
			)
		}
	}
}

// bottomAliens returns a map of aliens that are the lowest of its column
func (ag *alienGrid) bottomAliens() map[int]*alien {
	bottomAliens := make(map[int]*alien, ag.c.cols)
	for r := range ag.alienGridPos {
		for c, a := range ag.alienGridPos[r] {
			if ag.alienGridPos[r][c] != nil {
				bottomAliens[c] = a
			}
		}
	}

	return bottomAliens
}
