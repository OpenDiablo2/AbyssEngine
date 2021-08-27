package common

import lua "github.com/yuin/gopher-lua"

// LuaTypeExport is a collection of all the information we need about a Go type in order to export it for
// use in Lua scripts.
type LuaTypeExport struct {
	Name            string
	ConstructorFunc lua.LGFunction
	Methods         map[string]lua.LGFunction
}
