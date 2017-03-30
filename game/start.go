package game

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

// start holds the start screen state
type start struct {
	r            *sdl.Renderer
	t1           *sdl.Texture
	t2           *sdl.Texture
	tx           int32
	ty           int32
	tw           int32
	th           int32
	titleFont    *ttf.Font
	infoFont     *ttf.Font
	frameCounter int
}

// newStart returns a new start screen
func newStart(r *sdl.Renderer) (*start, error) {
	maxX, maxY, _ := r.GetRendererOutputSize()
	s := &start{
		r:            r,
		tw:           400,
		th:           428,
		frameCounter: 0,
	}

	// Set texture
	var err error
	s.t1, err = img.LoadTexture(r, "assets/alien_l1.png")
	if err != nil {
		return nil, fmt.Errorf("couldn't create start texture 1: %v", err)
	}
	s.t2, err = img.LoadTexture(r, "assets/alien_l2.png")
	if err != nil {
		return nil, fmt.Errorf("couldn't create start texture 2: %v", err)
	}

	// Set position
	s.tx = int32(maxX)/2 - s.tw/2
	s.ty = int32(maxY)/2 - s.th/2 - 100

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
	s.frameCounter++
	if s.frameCounter < 10 {
		s.r.Copy(s.t1, nil, &sdl.Rect{X: s.tx, Y: s.ty, W: s.tw, H: s.th})
	} else {
		s.r.Copy(s.t2, nil, &sdl.Rect{X: s.tx, Y: s.ty, W: s.tw, H: s.th})
	}
	if s.frameCounter > 20 {
		s.frameCounter = 0
	}

	maxX, maxY, _ := s.r.GetRendererOutputSize()

	title, _ := s.titleFont.RenderUTF8_Solid(
		"lil' e invaders",
		sdl.Color{R: 0xF6, G: 0x25, B: 0x9B, A: 0},
	)
	defer title.Free()

	info, _ := s.infoFont.RenderUTF8_Solid(
		"PRESS ENTER TO START",
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
