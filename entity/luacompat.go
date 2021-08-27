package entity

import (
	"fmt"

	"github.com/OpenDiablo2/AbyssEngine/common"
	lua "github.com/yuin/gopher-lua"
)

var luaTypeExportName = "node"
var LuaTypeExport = common.LuaTypeExport{
	Name: luaTypeExportName,
	//ConstructorFunc: newLuaEntity,
	Methods: map[string]lua.LGFunction{
		"appendChild": luaAppendChild,
	},
}

func (e *Entity) ToLua(l *lua.LState) *lua.LUserData {
	result := l.NewUserData()
	result.Value = e

	l.SetMetatable(result, l.GetTypeMetatable(luaTypeExportName))

	return result
}

func FromLua(ud *lua.LUserData) (*Entity, error) {
	v, ok := ud.Value.(*Entity)

	if !ok {
		return nil, fmt.Errorf("failed to convert")
	}

	return v, nil
}

func luaAppendChild(l *lua.LState) int {
	if l.GetTop() != 2 {
		l.ArgError(1, "argument expected")
		return 0
	}

	self, ok := l.ToUserData(1).Value.(*Entity)

	if !ok {
		l.RaiseError("failed to convert")
		return 0
	}

	child, ok := l.ToUserData(2).Value.(*Entity)

	if !ok {
		l.RaiseError("failed to convert")
		return 0
	}

	self.AddChild(child)

	return 0
}
