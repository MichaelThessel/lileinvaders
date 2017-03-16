package game

import (
	"github.com/veandco/go-sdl2/sdl"
)

type player struct {
	r        *sdl.Renderer
	x        int32
	y        int32
	w        int32
	h        int32
	stepSize int32
}

func newPlayer(r *sdl.Renderer) *player {
	maxX, maxY, _ := r.GetRendererOutputSize()
	p := &player{
		r:        r,
		w:        50,
		h:        50,
		stepSize: 10,
	}

	p.x = int32(maxX)/2 - p.w/2
	p.y = int32(maxY) - p.h

	return p
}

func (p *player) Draw() {
	p.r.SetDrawColor(0xFF, 0, 0, 0xFF)

	p.r.FillRect(
		&sdl.Rect{X: p.x, Y: p.y, W: p.w, H: p.h},
	)
}

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
