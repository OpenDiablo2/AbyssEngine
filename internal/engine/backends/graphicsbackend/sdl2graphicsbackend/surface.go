package sdl2graphicsbackend

import (
	"github.com/OpenDiablo2/AbyssEngine/internal/engine/backends/graphicsbackend"
	"github.com/veandco/go-sdl2/sdl"
)

var _ graphicsbackend.Surface = &SDL2Surface{}

type SDL2Surface struct {
	texture       *sdl.Texture
	width, height int32
	renderer      *sdl.Renderer
}

func CreateSDL2Surface(renderer *sdl.Renderer, 	texture *sdl.Texture, width, height int32) *SDL2Surface {
	result := &SDL2Surface{
		renderer: renderer,
		texture:  texture,
		width:    width,
		height:   height,
	}

	return result
}

func (s SDL2Surface) RenderTo(targetSurface graphicsbackend.Surface) error {
	sdlTargetSurface := targetSurface.(*SDL2Surface)

	destRect := &sdl.Rect{
		W: s.width, H: s.height,
	}

	if err := s.renderer.SetRenderTarget(sdlTargetSurface.texture); err != nil {
		return err
	}

	if err := s.renderer.Copy(s.texture, nil, destRect); err != nil {
		return err
	}

	if err := s.renderer.SetRenderTarget(nil); err != nil {
		return err
	}

	return nil
}
