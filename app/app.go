package app

import "github.com/veandco/go-sdl2/sdl"

type Config struct {
	Width     int
	Height    int
	Title     string
	FrameRate uint32
}

type App struct {
	w    *sdl.Window
	r    *sdl.Renderer
	c    *Config
	quit chan bool
}

func New(c *Config) (*App, error) {
	a := &App{c: c}
	if err := a.setup(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) setup() error {
	if err := a.setupWindow(); err != nil {
		return err
	}

	if err := a.setupRenderer(); err != nil {
		return err
	}

	return nil
}

func (a *App) Run() int {
	a.quit = make(chan bool)

loop:
	for {
		a.handleEvents()

		select {
		case <-a.quit:
			break loop
		default:
		}

		a.setBackground()

		a.r.Present()
		sdl.Delay(1000 / a.c.FrameRate)
	}

	return 0
}

func (a *App) setupWindow() error {
	var err error

	sdl.Do(func() {
		a.w, err = sdl.CreateWindow(
			a.c.Title,
			sdl.WINDOWPOS_UNDEFINED,
			sdl.WINDOWPOS_UNDEFINED,
			a.c.Width,
			a.c.Height,
			sdl.WINDOW_OPENGL,
		)
	})

	return err
}

func (a *App) setupRenderer() error {
	var err error
	sdl.Do(func() {
		a.r, err = sdl.CreateRenderer(a.w, -1, sdl.RENDERER_ACCELERATED)
	})

	if err != nil {
		return err
	}

	sdl.Do(func() {
		a.setBackground()
	})

	return nil
}

func (a *App) setBackground() {
	a.r.Clear()
	a.r.SetDrawColor(0, 0, 0, 0xFF)
	a.r.FillRect(
		&sdl.Rect{X: 0, Y: 0, W: int32(a.c.Width), H: int32(a.c.Height)},
	)
}

func (a *App) handleEvents() {
	sdl.Do(func() {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch e.(type) {
			case *sdl.QuitEvent:
				a.quit <- true
			case *sdl.KeyDownEvent:
				switch e.(*sdl.KeyDownEvent).Keysym.Sym {
				case sdl.K_q:
					go func() {
						a.quit <- true
						close(a.quit)
					}()
				}
			}
		}
	})
}

func (a *App) Destroy() {
	sdl.Do(func() {
		a.w.Destroy()
	})

	sdl.Do(func() {
		a.r.Destroy()
	})
}
