package game

import "github.com/veandco/go-sdl2/sdl"

// bullet holds bullet state information
type bullet struct {
	r         *sdl.Renderer
	x         int32
	y         int32
	w         int32
	h         int32
	direction int32
	speed     int32
}

// newBullet renerates a new bullet and adds it to the bullet list
func newBullet(r *sdl.Renderer, bl *bulletList, x, y, direction int32) {
	b := &bullet{
		r:         r,
		x:         x,
		y:         y,
		w:         3,
		h:         5,
		direction: direction, // -1 up 1 down
		speed:     25,
	}

	*bl = append(*bl, b)

	b.Draw()

}

// Draw an individual bullet
func (b *bullet) Draw() {
	b.r.SetDrawColor(0, 0xFF, 0, 0xFF)

	b.r.FillRect(
		&sdl.Rect{X: b.x, Y: b.y, W: b.w, H: b.h},
	)
}

// Update updates a bullets position
// This will return false if the bullet is out of bounds
func (b *bullet) Update() bool {
	_, maxY, _ := b.r.GetRendererOutputSize()

	b.y += b.direction * b.speed

	return !(b.y < 0 || b.y > int32(maxY))
}

// Holds all bullets currently on the screen
type bulletList []*bullet

// Draw renders all existing bullets
func (bl *bulletList) Draw() {
	tmpBl := bulletList{}
	for _, b := range *bl {
		if b.Update() {
			b.Draw()
			tmpBl = append(tmpBl, b)
		}
	}
	*bl = tmpBl
}
