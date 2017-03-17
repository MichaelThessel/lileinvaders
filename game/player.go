package game

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

// player holds the player state
type player struct {
	r        *sdl.Renderer
	t        *sdl.Texture
	x        int32
	y        int32
	w        int32
	h        int32
	stepSize int32
}

// newPlayer generates a player
func newPlayer(r *sdl.Renderer) (*player, error) {
	maxX, maxY, _ := r.GetRendererOutputSize()
	p := &player{
		r:        r,
		w:        90,
		h:        54,
		stepSize: 15,
	}

	var err error
	p.t, err = img.LoadTexture(r, "assets/tank.png")
	if err != nil {
		return nil, fmt.Errorf("couldn't create player texture: %v", err)
	}

	p.x = int32(maxX)/2 - p.w/2
	p.y = int32(maxY) - p.h

	return p, nil
}

// Draw draws the player
func (p *player) Draw() {
	p.r.Copy(p.t, nil, &sdl.Rect{X: p.x, Y: p.y, W: p.w, H: p.h})
}

// Move moves the player in a given direction
func (p *player) Move(direction rune) {
	maxX, _, _ := p.r.GetRendererOutputSize()
	switch direction {
	case 'l':
		p.x -= p.stepSize
		if p.x < 0 {
			p.x = 0
		}
	case 'r':
		p.x += p.stepSize
		if p.x+p.w > int32(maxX) {
			p.x = int32(maxX) - p.w
		}
	}
}

// Fire fires a bullet
func (p *player) Fire(bullets *bulletList) {
	if len(*bullets) < 1 {
		newBullet(p.r, bullets, p.x+p.w/2, p.y, -1)
	}
}
