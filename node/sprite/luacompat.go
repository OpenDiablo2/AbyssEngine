package sprite

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/OpenDiablo2/AbyssEngine/common"
	lua "github.com/yuin/gopher-lua"
)

var luaTypeExportName = "sprite"
var LuaTypeExport = common.LuaTypeExport{
	Name: luaTypeExportName,
	//ConstructorFunc: newLuaEntity,
	Methods: map[string]lua.LGFunction{
		"node":                   luaGetNode,
		"cellSize":               luaGetSetCellSize,
		"active":                 luaGetSetActive,
		"visible":                luaGetSetVisible,
		"position":               luaGetSetPosition,
		"currentSequence":        luaGetSetCurrentSequence,
		"currentFrame":           luaGetSetCurrentFrame,
		"sequenceCount":          luaGetSequenceCount,
		"frameCount":             luaGetFrameCount,
		"destroy":                luaDestroy,
		"mouseButtonDownHandler": luaGetSetMouseButtonDownHandler,
		"mouseButtonUpHandler":   luaGetSetMouseButtonUpHandler,
		"mouseOverHandler":       luaGetSetMouseOverHandler,
		"mouseLeaveHandler":      luaGetSetMouseLeaveHandler,
	},
}

func luaGetSetMouseOverHandler(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(l.NewFunction(func(l *lua.LState) int {
			sprite.onMouseOver()
			return 0
		}))

		return 1
	}

	luaFunc := l.CheckFunction(2)
	sprite.onMouseOver = func() {
		if err := l.CallByParam(lua.P{
			Fn:      luaFunc,
			NRet:    1,
			Protect: true,
		}, sprite.ToLua(l)); err != nil {
			panic(err)
		}
	}

	return 0
}

func luaGetSetMouseLeaveHandler(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(l.NewFunction(func(l *lua.LState) int {
			sprite.onMouseLeave()
			return 0
		}))

		return 1
	}

	luaFunc := l.CheckFunction(2)
	sprite.onMouseLeave = func() {
		if err := l.CallByParam(lua.P{
			Fn:      luaFunc,
			NRet:    1,
			Protect: true,
		}, sprite.ToLua(l)); err != nil {
			panic(err)
		}
	}

	return 0
}

func luaDestroy(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	sprite.Destroy()

	return 0
}

func luaGetSetMouseButtonUpHandler(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(l.NewFunction(func(l *lua.LState) int {
			sprite.onMouseButtonUp()
			return 0
		}))

		return 1
	}

	luaFunc := l.CheckFunction(2)
	sprite.onMouseButtonUp = func() {
		if err := l.CallByParam(lua.P{
			Fn:      luaFunc,
			NRet:    1,
			Protect: true,
		}, sprite.ToLua(l)); err != nil {
			panic(err)
		}
	}

	return 0
}

func luaGetFrameCount(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	l.Push(lua.LNumber(sprite.Sequences.FrameCount(sprite.CurrentSequence)))
	return 1
}

func luaGetSequenceCount(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	l.Push(lua.LNumber(sprite.Sequences.SequenceCount()))
	return 1
}

func luaGetSetCurrentFrame(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(sprite.CurrentFrame))
		return 1
	}

	newFrame := l.CheckInt(2)

	if (newFrame < 0) || (newFrame >= sprite.Sequences.FrameCount(sprite.CurrentSequence)) {
		l.RaiseError("frame index out of bounds")
		return 0
	}

	sprite.CurrentFrame = newFrame
	sprite.initialized = false

	return 0
}

func luaGetSetCurrentSequence(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(sprite.CurrentSequence))
		return 1
	}

	newSequence := l.CheckInt(2)

	if (newSequence < 0) || (newSequence >= sprite.Sequences.SequenceCount()) {
		l.RaiseError("sequence index out of bounds")
		return 0
	}

	sprite.CurrentSequence = newSequence
	sprite.CurrentFrame = 0
	sprite.initialized = false

	return 0
}

func luaGetNode(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	l.Push(sprite.Node.ToLua(l))

	return 1
}

func (s *Sprite) ToLua(l *lua.LState) *lua.LUserData {
	result := l.NewUserData()
	result.Value = s

	l.SetMetatable(result, l.GetTypeMetatable(luaTypeExportName))

	return result
}

func FromLua(ud *lua.LUserData) (*Sprite, error) {
	v, ok := ud.Value.(*Sprite)

	if !ok {
		return nil, fmt.Errorf("failed to convert")
	}

	return v, nil
}

func luaGetSetPosition(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(sprite.X))
		l.Push(lua.LNumber(sprite.Y))
		return 2
	}

	posX := l.ToNumber(2)
	posY := l.ToNumber(3)

	sprite.X = int(posX)
	sprite.Y = int(posY)

	return 0
}

func luaGetSetCellSize(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(sprite.CellSizeX))
		l.Push(lua.LNumber(sprite.CellSizeY))
		return 2
	}

	sizeX := l.ToNumber(2)
	sizeY := l.ToNumber(3)

	sprite.CellSizeX = int(sizeX)
	sprite.CellSizeY = int(sizeY)
	sprite.initialized = false
	rl.UnloadTexture(sprite.texture)

	return 0
}

func luaGetSetMouseButtonDownHandler(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(l.NewFunction(func(l *lua.LState) int {
			sprite.onMouseButtonDown()
			return 0
		}))

		return 1
	}

	luaFunc := l.CheckFunction(2)
	sprite.onMouseButtonDown = func() {
		if err := l.CallByParam(lua.P{
			Fn:      luaFunc,
			NRet:    1,
			Protect: true,
		}, sprite.ToLua(l)); err != nil {
			panic(err)
		}
	}

	return 0
}

func luaGetSetVisible(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LBool(sprite.Visible))
		return 1
	}

	newValue := l.CheckBool(2)
	sprite.Visible = newValue

	return 0
}

func luaGetSetActive(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LBool(sprite.Active))
		return 1
	}

	newValue := l.CheckBool(2)
	sprite.Active = newValue

	return 0
}
