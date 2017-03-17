package game

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

// alien holds the alien state
type alien struct {
	r        *sdl.Renderer
	t        *sdl.Texture
	x        int32
	y        int32
	w        int32
	h        int32
	stepSize int32
}

// newAlien generates a alien
func newAlien(r *sdl.Renderer, x, y int32) (*alien, error) {
	a := &alien{
		r:        r,
		w:        80,
		h:        86,
		x:        x,
		y:        y,
		stepSize: 10,
	}

	var err error
	a.t, err = img.LoadTexture(r, "assets/alien.png")
	if err != nil {
		return nil, fmt.Errorf("couldn't create alien texture: %v", err)
	}

	return a, nil
}

// Draw draws the alien
func (a *alien) Draw() {
	a.r.Copy(a.t, nil, &sdl.Rect{X: a.x, Y: a.y, W: a.w, H: a.h})
}

// Fire fires a bullet
func (a *alien) Fire(bullets *bulletList) {
	newBullet(a.r, bullets, a.x+a.w/2, a.y, -1)
}

// alienGrid holds the alien grid state
type alienGrid struct {
	r           *sdl.Renderer
	alienList   []*alien
	rows        int   // number of rows
	cols        int   // number of columns
	marginRow   int   // space between rows
	marginCol   int   // space between columns
	direction   int32 // direction of x movement (1: left, -1: right)
	returnPoint int32 // when to switch the x direction
	speed       int32 // grid movement speed
	speedStep   int   // after how many drops to increase the speed
	dropCount   int   // How often the grid moved down in y

}

// newAlienGrid creates a new alien grid
func newAlienGrid(renderer *sdl.Renderer) (*alienGrid, error) {
	maxX, _, _ := renderer.GetRendererOutputSize()

	ag := &alienGrid{
		r:           renderer,
		rows:        5,
		cols:        10,
		marginRow:   20,
		marginCol:   20,
		direction:   1,
		returnPoint: 30,
		speed:       4,
		speedStep:   5,
		dropCount:   0,
	}

	textureWidth := 80 //TODO: get this dynamically
	textureHeight := 86

	gridWidth := (textureWidth+ag.marginCol)*ag.cols - ag.marginCol

	startX := (maxX - gridWidth) / 2
	startY := 50

	currentX := startX
	currentY := startY
	for r := 0; r < ag.rows; r++ {
		for c := 0; c < ag.cols; c++ {
			a, err := newAlien(renderer, int32(currentX), int32(currentY))
			currentX += textureWidth + ag.marginCol
			if err != nil {
				return nil, err
			}
			ag.alienList = append(ag.alienList, a)
		}
		currentX = startX
		currentY += textureHeight + ag.marginRow
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
	// Viewport && grid dimensions
	maxX, _, _ := ag.r.GetRendererOutputSize()
	x1, _, x2, _ := ag.getDimensions()

	// Check if the grid hits the boundary
	moveY := false
	if x2 >= int32(maxX)-ag.returnPoint {
		ag.direction = -1
		moveY = true
	} else if x1 <= ag.returnPoint {
		ag.direction = 1
		moveY = true
	}

	// Increase the speed over time
	if moveY {
		ag.dropCount++
		if ag.dropCount%ag.speedStep == 0 {
			ag.speed++
		}
	}

	// Move all aliens
	for _, a := range ag.alienList {
		a.x += ag.direction * ag.speed
		if moveY {
			a.y += 3
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

// test hit checks if a bullet has hit an alien in the grid
func (ag *alienGrid) testHit(bl *bulletList) {
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
			break
		}
	}
}

// remove removes an alien from the grid
func (ag *alienGrid) remove(a *alien) {
	tmpAl := []*alien{}
	for _, ta := range ag.alienList {
		if ta != a {
			tmpAl = append(tmpAl, ta)
		}
	}

	ag.alienList = tmpAl
}
