package game

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

// start holds the start screen state
type start struct {
	r         *sdl.Renderer
	t         *sdl.Texture
	titleFont *ttf.Font
	infoFont  *ttf.Font
	x         int32
	y         int32
	w         int32
	h         int32
}

// newStart returns a new start screen
func newStart(r *sdl.Renderer) (*start, error) {
	maxX, maxY, _ := r.GetRendererOutputSize()
	s := &start{
		r: r,
		w: 400,
		h: 428,
	}

	// Set texture
	var err error
	s.t, err = img.LoadTexture(r, "assets/alien_l.png")
	if err != nil {
		return nil, fmt.Errorf("couldn't create start texture: %v", err)
	}

	// Set position
	s.x = int32(maxX)/2 - s.w/2
	s.y = int32(maxY)/2 - s.h/2 - 100

	// Set title font
	s.titleFont, err = ttf.OpenFont("assets/font.ttf", 80)
	if err != nil {
		return nil, fmt.Errorf("could not load font: %v", err)
	}

	// Set info font
	s.infoFont, err = ttf.OpenFont("assets/font.ttf", 20)
	if err != nil {
		return nil, fmt.Errorf("could not load font: %v", err)
	}

	return s, nil
}

// Draw draws the start screen
func (s *start) Draw() {
	s.r.Copy(s.t, nil, &sdl.Rect{X: s.x, Y: s.y, W: s.w, H: s.h})

	maxX, maxY, _ := s.r.GetRendererOutputSize()

	title, _ := s.titleFont.RenderUTF8_Solid(
		"lil' e invaders",
		sdl.Color{R: 0xF6, G: 0x25, B: 0x9B, A: 0},
	)
	defer title.Free()

	info, _ := s.infoFont.RenderUTF8_Solid(
		"Press Space to Start",
		sdl.Color{R: 0xF6, G: 0x25, B: 0x9B, A: 0},
	)
	defer info.Free()

	var clipRect sdl.Rect
	title.GetClipRect(&clipRect)
	info.GetClipRect(&clipRect)

	titleTex, _ := s.r.CreateTextureFromSurface(title)
	infoTex, _ := s.r.CreateTextureFromSurface(info)

	s.r.Copy(
		titleTex,
		nil,
		&sdl.Rect{
			X: int32(maxX)/2 - title.W/2,
			Y: int32(maxY) - 250,
			W: title.W,
			H: title.H,
		},
	)

	s.r.Copy(
		infoTex,
		nil,
		&sdl.Rect{
			X: int32(maxX)/2 - info.W/2,
			Y: int32(maxY) - 120,
			W: info.W,
			H: info.H,
		},
	)
}
