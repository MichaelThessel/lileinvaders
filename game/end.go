package game

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

// end holds the end screen state
type end struct {
	r         *sdl.Renderer
	scoreFont *ttf.Font
	infoFont  *ttf.Font
	score     int
}

// newEnd returns a new end screen
func newEnd(r *sdl.Renderer, score int) (*end, error) {
	e := &end{
		r:     r,
		score: score,
	}

	var err error

	// Set score font
	e.scoreFont, err = ttf.OpenFont("assets/font.ttf", 80)
	if err != nil {
		return nil, fmt.Errorf("could not load font: %v", err)
	}

	// Set info font
	e.infoFont, err = ttf.OpenFont("assets/font.ttf", 20)
	if err != nil {
		return nil, fmt.Errorf("could not load font: %v", err)
	}

	return e, nil
}

// Draw draws the end screen
func (e *end) Draw() {
	maxX, maxY, _ := e.r.GetRendererOutputSize()

	score, _ := e.scoreFont.RenderUTF8_Solid(
		fmt.Sprintf("POINTS: %d", e.score),
		sdl.Color{R: 0xF6, G: 0x25, B: 0x9B, A: 0},
	)
	defer score.Free()

	info1, _ := e.scoreFont.RenderUTF8_Solid(
		"GAME OVER",
		sdl.Color{R: 0xF6, G: 0x25, B: 0x9B, A: 0},
	)
	defer info1.Free()

	info2, _ := e.infoFont.RenderUTF8_Solid(
		"PRESS ENTER TO RESTART",
		sdl.Color{R: 0xF6, G: 0x25, B: 0x9B, A: 0},
	)
	defer info2.Free()

	var clipRect sdl.Rect
	score.GetClipRect(&clipRect)
	info1.GetClipRect(&clipRect)
	info2.GetClipRect(&clipRect)

	scoreTex, _ := e.r.CreateTextureFromSurface(score)
	info1Tex, _ := e.r.CreateTextureFromSurface(info1)
	info2Tex, _ := e.r.CreateTextureFromSurface(info2)

	e.r.Copy(
		scoreTex,
		nil,
		&sdl.Rect{
			X: int32(maxX)/2 - score.W/2,
			Y: int32(maxY) / 2,
			W: score.W,
			H: score.H,
		},
	)

	e.r.Copy(
		info1Tex,
		nil,
		&sdl.Rect{
			X: int32(maxX)/2 - info1.W/2,
			Y: int32(maxY)/2 - 100,
			W: info1.W,
			H: info1.H,
		},
	)

	e.r.Copy(
		info2Tex,
		nil,
		&sdl.Rect{
			X: int32(maxX)/2 - info2.W/2,
			Y: int32(maxY)/2 + 100,
			W: info2.W,
			H: info2.H,
		},
	)
}
