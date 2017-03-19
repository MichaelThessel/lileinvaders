package game

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

type stats struct {
	r      *sdl.Renderer
	font   *ttf.Font
	lifes  int
	points int
}

func newStats(r *sdl.Renderer, lifes int) (*stats, error) {
	s := &stats{
		r:      r,
		lifes:  lifes,
		points: 0,
	}

	var err error
	s.font, err = ttf.OpenFont("assets/font.ttf", 40)
	if err != nil {
		return nil, fmt.Errorf("could not load font: %v", err)
	}

	return s, nil
}

// Draw draws the stats
func (s *stats) Draw(lifes, points int) {
	maxX, _, _ := s.r.GetRendererOutputSize()

	lsf, _ := s.font.RenderUTF8_Solid(
		fmt.Sprintf("LIFES: %d", lifes),
		sdl.Color{R: 0xF6, G: 0x25, B: 0x9B, A: 0},
	)
	defer lsf.Free()

	psf, _ := s.font.RenderUTF8_Solid(
		fmt.Sprintf("POINTS: %08d", points),
		sdl.Color{R: 0xF6, G: 0x25, B: 0x9B, A: 0},
	)
	defer psf.Free()

	var clipRect sdl.Rect
	lsf.GetClipRect(&clipRect)
	psf.GetClipRect(&clipRect)

	lsfTex, _ := s.r.CreateTextureFromSurface(lsf)
	psfTex, _ := s.r.CreateTextureFromSurface(psf)

	s.r.Copy(lsfTex, nil, &sdl.Rect{X: 10, Y: 10, W: lsf.W, H: lsf.H})
	s.r.Copy(psfTex, nil, &sdl.Rect{X: int32(maxX) - psf.W - 10, Y: 10, W: psf.W, H: psf.H})
}
