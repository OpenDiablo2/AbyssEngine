package sdl2graphicsbackend

import (
	"github.com/OpenDiablo2/AbyssEngine/internal/engine/backends/graphicsbackend"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

var _ graphicsbackend.Interface = &SDL2GraphicsBackend{}

type SDL2GraphicsBackend struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

func (r *SDL2GraphicsBackend) Render() error {
	if err := r.renderer.Clear(); err != nil {
		return err
	}

	r.renderer.Present()

	return nil
}

func Create() (*SDL2GraphicsBackend, error) {
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS); err != nil {
		return nil, err
	}

	window, err := sdl.CreateWindow("OpenDiablo 2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 640, 480, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE|sdl.WINDOW_INPUT_FOCUS)
	if err != nil {
		return nil, err
	}

	r, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_TARGETTEXTURE|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return nil, err
	}

	r.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	r.SetIntegerScale(false)
	r.SetLogicalSize(800, 600)
	window.SetMinimumSize(800, 600)

	result := &SDL2GraphicsBackend{
		window:   window,
		renderer: r,
	}

	return result, nil
}

func (r *SDL2GraphicsBackend) GetRendererName() string {
	return "SDL2"
}

func (r *SDL2GraphicsBackend) SetWindowIcon(fileName string) {
	panic("implement me")
}

func (r *SDL2GraphicsBackend) IsFullScreen() bool {
	panic("implement me")
}

func (r *SDL2GraphicsBackend) SetFullScreen(fullScreen bool) {
	panic("implement me")
}

func (r *SDL2GraphicsBackend) SetVSyncEnabled(vsync bool) {
	panic("implement me")
}

func (r *SDL2GraphicsBackend) GetVSyncEnabled() bool {
	panic("implement me")
}

func (r *SDL2GraphicsBackend) GetCursorPos() (int, int) {
	panic("implement me")
}

func (r *SDL2GraphicsBackend) CurrentFPS() float64 {
	panic("implement me")
}
