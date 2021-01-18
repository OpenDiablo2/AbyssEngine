package sdl2inputbackend

import (
	"os"

	"github.com/OpenDiablo2/AbyssEngine/internal/engine/backends/inputbackend"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	mousePressThreshold = 25
)

type mouseEventInfo struct {
	state bool
	time  uint32
}

var _ inputbackend.Interface = &SDL2InputBackend{}

type SDL2InputBackend struct {
	mouseState [8]mouseEventInfo
	cursorPosX int
	cursorPosY int
}

func Create() (*SDL2InputBackend, error) {
	result := &SDL2InputBackend{}

	return result, nil
}

func (i SDL2InputBackend) Process() error {
	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			os.Exit(0)
		case *sdl.MouseMotionEvent:
			i.cursorPosX = int(t.X)
			i.cursorPosY = int(t.Y)
		case *sdl.MouseButtonEvent:
			i.mouseState[t.Button].time = sdl.GetTicks()
			i.mouseState[t.Button].state = t.Type == sdl.MOUSEBUTTONDOWN
		}
	}

	return nil
}

func (i *SDL2InputBackend) CursorPosition() (x int, y int) {
	return i.cursorPosX, i.cursorPosY
}

func (i *SDL2InputBackend) InputChars() []rune {
	return []rune{}
}

func (i *SDL2InputBackend) IsKeyPressed(key inputbackend.Key) bool {
	return false
}

func (i *SDL2InputBackend) IsKeyJustPressed(key inputbackend.Key) bool {
	return false
}

func (i *SDL2InputBackend) IsKeyJustReleased(key inputbackend.Key) bool {
	return false
}

func (i *SDL2InputBackend) IsMouseButtonPressed(button inputbackend.MouseButton) bool {
	switch button {
	case inputbackend.MouseButtonLeft:
		return i.mouseState[1].state
	case inputbackend.MouseButtonRight:
		return i.mouseState[2].state
	case inputbackend.MouseButtonMiddle:
		return i.mouseState[3].state
	default:
		return false
	}
}

func (i *SDL2InputBackend) IsMouseButtonJustPressed(button inputbackend.MouseButton) bool {
	switch button {
	case inputbackend.MouseButtonLeft:
		return i.mouseState[1].state == true && (sdl.GetTicks()-i.mouseState[1].time) < mousePressThreshold
	case inputbackend.MouseButtonRight:
		return i.mouseState[2].state == true && (sdl.GetTicks()-i.mouseState[2].time) < mousePressThreshold
	case inputbackend.MouseButtonMiddle:
		return i.mouseState[3].state == true && (sdl.GetTicks()-i.mouseState[3].time) < mousePressThreshold
	default:
		return false
	}
}

func (i *SDL2InputBackend) IsMouseButtonJustReleased(button inputbackend.MouseButton) bool {
	switch button {
	case inputbackend.MouseButtonLeft:
		return i.mouseState[1].state == false && (sdl.GetTicks()-i.mouseState[1].time) < mousePressThreshold
	case inputbackend.MouseButtonRight:
		return i.mouseState[2].state == false && (sdl.GetTicks()-i.mouseState[2].time) < mousePressThreshold
	case inputbackend.MouseButtonMiddle:
		return i.mouseState[3].state == false && (sdl.GetTicks()-i.mouseState[3].time) < mousePressThreshold
	default:
		return false
	}
}

func (i *SDL2InputBackend) KeyPressDuration(key inputbackend.Key) int {
	return 0
}
