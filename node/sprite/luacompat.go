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
		"mouseButtonDownHandler": luaGetSetMouseButtonDownHandler,
	},
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

func (e *Sprite) ToLua(l *lua.LState) *lua.LUserData {
	result := l.NewUserData()
	result.Value = e

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

//
//func (s *Sprite) IndexSet(index, value tengo.Object) error {
//	indexStr, ok := tengo.ToString(index)
//
//	if !ok {
//		return errors.New("invalid index")
//	}
//
//	switch indexStr {
//	case "onMouseButtonDown":
//		fn, ok := tengo.ToInterface(value).(*tengo.CompiledFunction)
//		if !ok {
//			s.onMouseButtonDown = nil
//		} else {
//			s.onMouseButtonDown = fn
//		}
//
//		return nil
//	case "onMouseButtonUp":
//		fn, ok := tengo.ToInterface(value).(*tengo.CompiledFunction)
//		if !ok {
//			s.onMouseButtonUp = nil
//		} else {
//			s.onMouseButtonUp = fn
//		}
//
//		return nil
//	}
//
//	return fmt.Errorf("invalid index: %s", indexStr)
//}
