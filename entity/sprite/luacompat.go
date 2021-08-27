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
		"node":        luaGetNode,
		"setCellSize": luaSetCellSize,
		"active":      luaGetSetActive,
		"visible":     luaGetSetVisible,
	},
}

func luaGetNode(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	l.Push(sprite.Entity.ToLua(l))

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

func luaSetCellSize(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	sizeX := l.ToNumber(2)
	sizeY := l.ToNumber(3)

	sprite.CellSizeX = int(sizeX)
	sprite.CellSizeY = int(sizeY)
	sprite.initialized = false
	rl.UnloadTexture(sprite.texture)

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

//func (s *Sprite) IndexGet(index tengo.Object) (value tengo.Object, err error) {
//	indexStr, ok := tengo.ToString(index)
//
//	if !ok {
//		return nil, errors.New("invalid index")
//	}
//
//	switch indexStr {
//	case "x":
//		return &tengo.Int{Value: int64(s.X)}, nil
//	case "y":
//		return &tengo.Int{Value: int64(s.Y)}, nil
//	case "node":
//		return s.Entity, nil
//	case "setPosition":
//		return &tengo.UserFunction{
//			Name: "appendChild",
//			Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
//				if len(args) != 2 {
//					return nil, errors.New("expected two arguments")
//				}
//
//				posX, ok := tengo.ToInt(args[0])
//
//				if !ok {
//					return nil, errors.New("first argument must be int")
//				}
//
//				posY, ok := tengo.ToInt(args[0])
//
//				if !ok {
//					return nil, errors.New("first argument must be int")
//				}
//
//				s.X = posX
//				s.Y = posY
//
//				return s, nil
//			},
//		}, nil
//	}
//
//	return nil, fmt.Errorf("invalid index: %s", indexStr)
//
//}
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
