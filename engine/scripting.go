package engine

import (
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/OpenDiablo2/AbyssEngine/common"
	"github.com/OpenDiablo2/AbyssEngine/entity"
	"github.com/OpenDiablo2/AbyssEngine/entity/sprite"
	"github.com/OpenDiablo2/AbyssEngine/loader/filesystemloader"
	"github.com/OpenDiablo2/AbyssEngine/loader/mpqloader"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	lua "github.com/yuin/gopher-lua"
)

var luaTypes = []common.LuaTypeExport{
	entity.LuaTypeExport,
	sprite.LuaTypeExport,
}

func registerType(l *lua.LState, luaTypeExport common.LuaTypeExport) {
	typeMetatable := l.NewTypeMetatable(luaTypeExport.Name)
	l.SetGlobal(luaTypeExport.Name, typeMetatable)

	// static attributes
	if luaTypeExport.ConstructorFunc != nil {
		l.SetField(typeMetatable, "new", l.NewFunction(luaTypeExport.ConstructorFunc))
	}

	// methods
	l.SetField(typeMetatable, "__index", l.SetFuncs(l.NewTable(), luaTypeExport.Methods))
}

func (e *Engine) bootstrapScripts() {
	go func() {
		time.Sleep(time.Second * 1)

		e.luaState = lua.NewState()

		for _, luaType := range luaTypes {
			registerType(e.luaState, luaType)
		}

		// Inject searcher for loading files
		lPackage := e.luaState.GetGlobal("package")
		lSearchers := e.luaState.GetField(lPackage, "loaders").(*lua.LTable)

		// TODO: Try to remove normal searchers or something for security
		//for lSearchers.Len() > 0 {
		//	lSearchers.Remove(0)
		//}

		lSearchers.Append(e.luaState.NewFunction(func(l *lua.LState) int { return e.luaLoader(l) }))

		e.luaState.PreloadModule("abyss", func(l *lua.LState) int {
			mod := l.SetFuncs(e.luaState.NewTable(), map[string]lua.LGFunction{
				// shutdown()
				// shuts down the engine (exits to desktop)
				"shutdown": func(l *lua.LState) int { return e.luaShutdown(l) },

				// log(level: string, message: string)
				// writes a log entry message based on the log level: info, error, fatal, warn, debug, trace
				"log": func(l *lua.LState) int { return e.luaLog(l) },

				// fmt(format: string, values: any...)
				// returns a formatted string from the input format and values
				"fmt": func(l *lua.LState) int { return e.luaFmt(l) },

				// setBootText(text: string)
				// sets the lower text in the boot splash screen
				"setBootText": func(l *lua.LState) int { return e.luaSetBootText(l) },

				// joinPath(path: string...)
				// Joins path names together in an os-specific way
				"joinPath": func(l *lua.LState) int { return e.luaJoinPath(l) },

				// sleep(msec: int)
				// sleeps the specified number of milliseconds
				"sleep": func(l *lua.LState) int { return e.luaSleep(l) },

				// getEngineSettings()
				// returns the engine settings
				"getEngineSettings": func(l *lua.LState) int { return e.luaGetEngineSettings(l) },

				// addLoaderProvider(type: string, path: string)
				// adds a loader to the engine
				"addLoaderProvider": func(l *lua.LState) int { return e.luaAddLoaderProvider(l) },

				// exitBootMode()
				// exits boot mode and starts the main rendering system
				"exitBootMode": func(l *lua.LState) int { return e.luaExitBootMode(l) },

				// loadString(path: string)
				// loads a string from the loader
				"loadString": func(l *lua.LState) int { return e.luaExitBootMode(l) },

				// splitString(source: string, splitChars: string)
				// splits a string by the specified split characters
				"luaSplitString": func(l *lua.LState) int { return e.luaSplitString(l) },

				// getRootNode()
				// returns the root entity node for the engine
				"getRootNode": func(l *lua.LState) int { return e.luaGetRootNode(l) },

				// loadSprite(filePath: string, palette: string) Sprite
				// returns a sprite based on the path and palette
				"loadSprite": func(l *lua.LState) int { return e.luaLoadSprite(l) },

				// setCursor(cursor: Sprite)
				// sets the current cursor, or clears it if nil
				"setCursor": func(l *lua.LState) int { return e.luaSetCursor(l) },

				// loadPalette(name: string, filePath: string)
				// loads palette for use with sprites
				"loadPalette": func(l *lua.LState) int { return e.luaLoadPalette(l) },
			})

			e.luaState.Push(mod)

			return 1
		})

		err := e.luaState.DoString("require(\"/bootstrap\")")

		if err != nil {
			e.panic(err.Error())
			return
		}

	}()
}

func (e *Engine) luaLoader(l *lua.LState) int {
	toLoad := l.CheckString(1)

	if len(toLoad) == 0 {
		l.Push(lua.LString(fmt.Sprintf("no resource name specified", toLoad)))
		return 1
	}

	toLoad = toLoad + ".lua"

	if !strings.HasPrefix(toLoad, "/") {
		toLoad = "/" + toLoad
	}

	log.Trace().Msgf("Loading file: %s", toLoad)

	bootScriptFile, err := e.loader.Load(toLoad)

	if err != nil {
		l.Push(lua.LString(fmt.Sprintf("%s: unknown resource", toLoad)))
		return 1
	}

	resource, err := l.Load(bootScriptFile, toLoad)

	if err != nil {
		l.Push(lua.LString(fmt.Sprintf("error loading '%s': %s", toLoad, err.Error())))
		return 1
	}

	l.Push(resource)
	return 1
}

func (e *Engine) luaShutdown(l *lua.LState) int {
	if l.GetTop() > 0 {
		l.ArgError(1, "no arguments expected")
		return 0
	}

	e.shutdown = true
	log.Info().Msg("engine shutting down due to script SHUTDOWN command")
	return 0
}

func (e *Engine) luaFmt(l *lua.LState) int {
	if l.GetTop() < 2 {
		l.ArgError(l.GetTop(), "expected at least 2 arguments")
		return 0
	}

	f := l.CheckString(1)
	fmtArgs := make([]interface{}, l.GetTop()-1)

	for i := 1; i < l.GetTop(); i++ {
		argType := l.Get(i + 1).Type().String()
		switch argType {
		case "string":
			fmtArgs[i-1] = l.CheckString(i + 1)
		case "number":
			fmtArgs[i-1] = l.CheckNumber(i + 1)
		case "boolean":
			fmtArgs[i-1] = l.CheckBool(i + 1)
		default:
			l.ArgError(l.GetTop(), fmt.Sprintf("unsupported argument type: %s", argType))
			return 0
		}
	}

	l.Push(lua.LString(fmt.Sprintf(f, fmtArgs...)))
	return 1
}

func (e *Engine) luaSetBootText(l *lua.LState) int {
	if l.GetTop() != 1 {
		l.ArgError(l.GetTop(), "expected one argument")
		return 0
	}

	str := l.CheckString(1)

	e.bootLoadText = str

	return 0
}

func (e *Engine) luaExitBootMode(l *lua.LState) int {
	if l.GetTop() > 0 {
		l.ArgError(1, "no arguments expected")
		return 0
	}

	log.Info().Msg("Entering game mode")
	e.engineMode = EngineModeGame

	return 0
}

func (e *Engine) luaJoinPath(l *lua.LState) int {
	str := make([]string, l.GetTop())

	for i := 0; i < l.GetTop(); i++ {
		str[i] = l.CheckString(i + 1)
	}

	l.Push(lua.LString(path.Join(str...)))
	return 1
}

func getTableForObject(l *lua.LState, source interface{}) lua.LValue {
	result := l.NewTable()

	r := reflect.ValueOf(source)

	for idx := 0; idx < reflect.Indirect(r).NumField(); idx++ {
		f := reflect.Indirect(r).Field(idx)

		fieldName := r.Type().Field(idx).Name

		switch f.Interface().(type) {
		case int:
			result.RawSetString(fieldName, lua.LNumber(f.Int()))
		case bool:
			result.RawSetString(fieldName, lua.LBool(f.Bool()))
		case string:
			result.RawSetString(fieldName, lua.LString(f.String()))
		case []string:
			strings := f.Interface().([]string)
			results := l.NewTable()
			for i := range strings {
				results.RawSetInt(i, lua.LString(strings[i]))
			}
			result.RawSetString(fieldName, results)

		case interface{}:
			result.RawSetString(fieldName, getTableForObject(l, f.Interface()))
		}
	}

	return result
}

func (e *Engine) luaGetEngineSettings(l *lua.LState) int {
	if l.GetTop() > 0 {
		l.ArgError(1, "no arguments expected")
		return 0
	}

	result := getTableForObject(l, e.config)
	l.Push(result)
	return 1
}

func (e *Engine) luaSleep(l *lua.LState) int {
	if l.GetTop() != 1 {
		l.ArgError(l.GetTop(), "expected one argument")
	}

	val := l.CheckInt(1)

	time.Sleep(time.Millisecond * time.Duration(val))

	return 0
}

func (e *Engine) luaAddLoaderProvider(l *lua.LState) int {
	if l.GetTop() != 2 {
		l.ArgError(l.GetTop(), "expected two arguments")
		return 0
	}

	loaderType := l.CheckString(1)
	p := l.CheckString(2)

	switch loaderType {
	case "mpq":
		provider, err := mpqloader.New(p)
		if err != nil {
			l.RaiseError(err.Error())
			return 0
		}

		e.loader.AddProvider(provider)
	case "filesystem":
		provider := filesystemloader.New(p)
		e.loader.AddProvider(provider)
	default:
		l.RaiseError("unknown loader type: %s", loaderType)
		return 0
	}

	return 0
}

func (e *Engine) luaLoadString(l *lua.LState) int {
	if l.GetTop() != 1 {
		l.ArgError(l.GetTop(), "expected one argument")
		return 0
	}

	val := l.CheckString(1)

	file, err := e.loader.Load(val)

	if err != nil {
		l.RaiseError(err.Error())
		return 0
	}

	result, err := ioutil.ReadAll(file)
	if err != nil {
		l.RaiseError(err.Error())
		return 0
	}

	l.Push(lua.LString(result))
	return 1
}

func (e *Engine) luaSplitString(l *lua.LState) int {
	if l.GetTop() != 2 {
		l.ArgError(l.GetTop(), "expected two arguments")
		return 0
	}

	sourceString := l.CheckString(1)
	splitChars := l.CheckString(2)

	strings := strings.Split(sourceString, splitChars)
	resultArray := l.NewTable()

	for i := 0; i < len(strings); i++ {
		resultArray.Append(lua.LString(strings[i]))
	}

	l.Push(resultArray)
	return 1
}

func (e *Engine) luaGetRootNode(l *lua.LState) int {
	if l.GetTop() > 0 {
		l.ArgError(1, "no arguments expected")
		return 0
	}

	l.Push(e.rootNode.ToLua(l))
	return 1
}

func (e *Engine) luaLoadSprite(l *lua.LState) int {
	if l.GetTop() != 2 {
		l.ArgError(l.GetTop(), "expected two arguments")
		return 0
	}

	filePath := l.CheckString(1)
	palette := l.CheckString(2)

	result, err := sprite.New(e.loader, e, filePath, palette)

	if err != nil {
		l.RaiseError(err.Error())
		return 0
	}

	l.Push(result.ToLua(l))
	return 1
}

func (e *Engine) luaSetCursor(l *lua.LState) int {
	if l.GetTop() != 1 {
		l.ArgError(l.GetTop(), "expected one argument")
		return 0
	}

	if l.Get(1).Type() == lua.LTNil {
		e.cursorSprite = nil
		return 0
	}

	sprite, err := sprite.FromLua(l.ToUserData(1))

	if err != nil {
		e.cursorSprite = nil
		l.RaiseError(err.Error())
		return 0
	}

	e.cursorSprite = sprite

	return 0
}

func (e *Engine) luaLog(l *lua.LState) int {
	if l.GetTop() != 2 {
		l.ArgError(l.GetTop(), "expected two arguments")
		return 0
	}

	var logObject *zerolog.Event

	logType := l.CheckString(1)

	switch logType {
	case "info":
		logObject = log.Info()
	case "error":
		logObject = log.Error()
	case "fatal":
		logObject = log.Fatal()
	case "warn":
		logObject = log.Warn()
	case "debug":
		logObject = log.Debug()
	case "trace":
		logObject = log.Trace()
	default:
		l.ArgError(1, "unexpected log type")
		return 0
	}

	m := l.CheckString(2)
	logObject.Msg(m)

	return 0
}

func (e *Engine) luaLoadPalette(l *lua.LState) int {
	if l.GetTop() != 2 {
		l.ArgError(l.GetTop(), "expected two arguments")
		return 0
	}

	palette := l.CheckString(1)
	filePath := l.CheckString(2)

	err := e.loadPalette(palette, filePath)

	if err != nil {
		l.RaiseError(err.Error())
		return 0
	}

	return 0
}
