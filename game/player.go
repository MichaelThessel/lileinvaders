package game

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
	mix "github.com/veandco/go-sdl2/sdl_mixer"
)

// playerConfig holds the player configuration
type playerConfig struct {
	stepSize    int32
	bulletSpeed int32
	lifes       int
}

// player holds the player state
type player struct {
	c      *playerConfig
	r      *sdl.Renderer
	t      *sdl.Texture
	sounds map[string]*mix.Chunk
	x      int32
	y      int32
	w      int32
	h      int32
	lifes  int
}

// newPlayer generates a player
func newPlayer(r *sdl.Renderer, c *playerConfig) (*player, error) {
	maxX, maxY, _ := r.GetRendererOutputSize()
	p := &player{
		c:     c,
		r:     r,
		w:     90,
		h:     54,
		lifes: c.lifes,
	}

	// Set texture
	var err error
	p.t, err = img.LoadTexture(r, "assets/tank.png")
	if err != nil {
		return nil, fmt.Errorf("couldn't create player texture: %v", err)
	}

	// Set position
	p.x = int32(maxX)/2 - p.w/2
	p.y = int32(maxY) - p.h

	// Set sounds
	p.sounds = make(map[string]*mix.Chunk, 0)
	p.sounds["fire"], err = mix.LoadWAV("assets/sounds/fire.wav")
	if err != nil {
		return nil, fmt.Errorf("couldn't load sound: %v", err)
	}
	p.sounds["hit"], err = mix.LoadWAV("assets/sounds/playerhit.wav")
	if err != nil {
		return nil, fmt.Errorf("couldn't load sound: %v", err)
	}

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
		p.x -= p.c.stepSize
		if p.x < 0 {
			p.x = 0
		}
	case 'r':
		p.x += p.c.stepSize
		if p.x+p.w > int32(maxX) {
			p.x = int32(maxX) - p.w
		}
	}
}

// Fire fires a bullet
func (p *player) Fire(bullets *bulletList) {
	if len(*bullets) > 0 {
		return
	}

	newBullet(
		p.r,
		bullets,
		p.x+p.w/2,
		p.y,
		&bulletConfig{
			speed:     p.c.bulletSpeed,
			direction: -1,
			colorR:    0x00,
			colorG:    0xFC,
			colorB:    0xFF,
		},
	)

	p.sounds["fire"].Play(0, 0)
}

// test hit checks if a bullet has hit player
func (p *player) testHit(bl *bulletList) (dead bool) {
	for _, b := range *bl {
		// Continue if bullet is beyond player dimensions
		if b.y+b.h < p.y || b.x+b.w < p.x || b.x > p.x+p.w {
			continue
		}

		bl.remove(b)

		p.sounds["hit"].Play(0, 0)

		p.lifes--
		if p.lifes == 0 {
			dead = true
			return
		}
	}

	return
}
