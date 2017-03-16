package app

import (
	"sort"

	"github.com/veandco/go-sdl2/sdl"
)

// Config holds the application configuration
type Config struct {
	Width     int
	Height    int
	Title     string
	FrameRate uint32
}

// App is the main application
type App struct {
	w               *sdl.Window
	r               *sdl.Renderer
	c               *Config
	quit            chan bool
	keyCallbacks    []keyCallback
	renderCallbacks renderCallbacks
}

// New returns a new app instance
func New(c *Config) (*App, error) {
	a := &App{c: c}
	if err := a.setup(); err != nil {
		return nil, err
	}

	return a, nil
}

// setup sets up the app
func (a *App) setup() error {
	if err := a.setupWindow(); err != nil {
		return err
	}

	if err := a.setupRenderer(); err != nil {
		return err
	}

	return nil
}

// Run starts the main app loop
// a bool true on the quit channel will break the loop and quit the app
func (a *App) Run() int {
	a.quit = make(chan bool)

	sort.Sort(a.renderCallbacks)

loop:
	for {
		a.clearWindow()

		a.handleEvents()

		select {
		case <-a.quit:
			break loop
		default:
		}

		for _, rc := range a.renderCallbacks {
			rc.callback()
		}

		a.r.Present()
		sdl.Delay(1000 / a.c.FrameRate)
	}

	return 0
}

// setupWindow sets up the app window
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

// setupRenderer sets up the renderer
func (a *App) setupRenderer() error {
	var err error
	sdl.Do(func() {
		a.r, err = sdl.CreateRenderer(a.w, -1, sdl.RENDERER_ACCELERATED)
	})

	if err != nil {
		return err
	}

	sdl.Do(func() {
		a.clearWindow()
	})

	return nil
}

// clearWindow clears the window
func (a *App) clearWindow() {
	a.r.Clear()
	a.r.SetDrawColor(0, 0, 0, 0xFF)
	a.r.FillRect(
		&sdl.Rect{X: 0, Y: 0, W: int32(a.c.Width), H: int32(a.c.Height)},
	)
}

// handleEvents handles input events
func (a *App) handleEvents() {
	sdl.Do(func() {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch e.(type) {
			case *sdl.QuitEvent:
				go func() {
					a.quit <- true
					close(a.quit)
				}()
			case *sdl.KeyDownEvent:
				switch e.(*sdl.KeyDownEvent).Keysym.Sym {
				case sdl.K_q:
					go func() {
						a.quit <- true
						close(a.quit)
					}()
				default:
					// Externally registered handlers
					for _, kh := range a.keyCallbacks {
						if kh.key == e.(*sdl.KeyDownEvent).Keysym.Sym {
							kh.callback()
						}
					}
				}
			}
		}
	})
}

// keyCallback defines key and callback associations
type keyCallback struct {
	key      sdl.Keycode
	callback func()
}

// RegisterKeyHandler allows to register keyboard event callbacks
func (a *App) RegisterKeyCallback(key sdl.Keycode, callback func()) {
	a.keyCallbacks = append(a.keyCallbacks, keyCallback{
		key:      key,
		callback: callback,
	})
}

// renderCallback defines callbacks to inject into the render loop
type renderCallback struct {
	priority int
	callback func()
}

// renderCallbacks holds render callbacks and implements sort.Interface
type renderCallbacks []*renderCallback

func (rc renderCallbacks) Less(i, j int) bool { return rc[i].priority < rc[j].priority }
func (rc renderCallbacks) Len() int           { return len(rc) }
func (rc renderCallbacks) Swap(i, j int)      { rc[i], rc[j] = rc[j], rc[i] }

// RegisterRenderCallback registers a callback that will be called on each
// render cycle
func (a *App) RegisterRenderCallback(priority int, callback func()) {
	a.renderCallbacks = append(a.renderCallbacks, &renderCallback{
		priority: priority,
		callback: callback,
	})
}

// GetRenderer returns a renderer instance
func (a *App) GetRenderer() *sdl.Renderer {
	return a.r
}

// Destroy destroys the app
func (a *App) Destroy() {
	sdl.Do(func() {
		a.w.Destroy()
	})

	sdl.Do(func() {
		a.r.Destroy()
	})
}
