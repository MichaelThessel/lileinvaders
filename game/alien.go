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
		h:        59,
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

// Move moves the alien in a given direction
func (a *alien) Move() {
}

// Fire fires a bullet
func (a *alien) Fire(bullets *bulletList) {
	newBullet(a.r, bullets, a.x+a.w/2, a.y, -1)
}

// alienGrid holds the alien grid state
type alienGrid struct {
	r           *sdl.Renderer
	alienList   []*alien
	rows        int
	cols        int
	marginRow   int
	marginCol   int
	direction   int32
	returnPoint int32
	speed       int32
}

// newAlienGrid creates a new alien grid
func newAlienGrid(renderer *sdl.Renderer) (*alienGrid, error) {
	maxX, _, _ := renderer.GetRendererOutputSize()

	ag := &alienGrid{
		r:           renderer,
		rows:        5,
		cols:        6,
		marginRow:   10,
		marginCol:   10,
		direction:   1,
		returnPoint: 30,
		speed:       3,
	}

	textureWidth := 80 //TODO: get this dynamically
	textureHeight := 59

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
	// Grid return points
	maxX, _, _ := ag.r.GetRendererOutputSize()

	x1, _, x2, _ := ag.getDimensions()

	moveY := false
	if x2 >= int32(maxX)-ag.returnPoint {
		ag.direction = -1
		moveY = true
	} else if x1 <= ag.returnPoint {
		ag.direction = 1
		moveY = true
	}

	for _, a := range ag.alienList {
		a.x += ag.direction * ag.speed
		if moveY {
			a.y += ag.speed
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