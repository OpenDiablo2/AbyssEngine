package engine

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/OpenDiablo2/AbyssEngine/entity/sprite"
	"github.com/OpenDiablo2/AbyssEngine/loader/filesystemloader"
	"github.com/OpenDiablo2/AbyssEngine/loader/mpqloader"
	"github.com/d5/tengo/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (e *Engine) bootstrapScripts() {
	go func() {
		time.Sleep(time.Second * 1)

		bootScriptFile, err := e.loader.Load("/bootstrap.tengo")

		if err != nil {
			e.panic("Could not locate bootstrap script.")
			return
		}

		bootScript, err := ioutil.ReadAll(bootScriptFile)

		if err != nil {
			e.panic("Could not locate bootstrap script.")
			return
		}

		bootScriptFile.Close()

		script := tengo.NewScript(bootScript)
		script.EnableFileImport(true)
		_ = script.SetImportDir(e.config.RootPath)

		e.addScriptFunctions(script)

		compiled, err := script.Compile()

		if err != nil {
			e.panic(err.Error())
			return
		}

		if err := compiled.Run(); err != nil {
			e.panic(err.Error())
		}
	}()
}

func (e *Engine) addScriptFunctions(script *tengo.Script) {

	mod := tengo.NewModuleMap()
	mod.AddBuiltinModule("abyss", map[string]tengo.Object{
		// fmt(format: string, values: any...)
		// returns a formatted string from the input format and values
		"fmt": &tengo.UserFunction{Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) < 2 {
				return tengo.UndefinedValue, errors.New("expected at least 2 arguments")
			}

			f, _ := tengo.ToString(args[0])
			fmtArgs := make([]interface{}, len(args)-1)

			for i := 1; i < len(args); i++ {
				typeName := args[i].TypeName()
				switch typeName {
				case "string":
					fmtArgs[i-1], _ = tengo.ToString(args[i])
				case "int":
					fmtArgs[i-1], _ = tengo.ToInt(args[i])
				case "bool":
					fmtArgs[i-1], _ = tengo.ToBool(args[i])
				default:
					return tengo.UndefinedValue, fmt.Errorf("unknown type: %s", typeName)
				}
			}

			resultString := fmt.Sprintf(f, fmtArgs...)

			return &tengo.String{Value: resultString}, nil
		}},

		// setBootText(text: string)
		// sets the lower text in the boot splash screen
		"setBootText": &tengo.UserFunction{Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 1 {
				return tengo.UndefinedValue, errors.New("expected only one argument")
			}

			str, ok := tengo.ToString(args[0])

			if !ok {
				return tengo.UndefinedValue, errors.New("argument must be a string")
			}

			e.bootLoadText = str

			return tengo.UndefinedValue, nil
		}},

		// log(level: string, message: string)
		// writes a log entry message based on the log level: info, error, fatal, warn, debug, trace
		"log": &tengo.UserFunction{Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 2 {
				return tengo.UndefinedValue, errors.New("expected 2 arguments")
			}

			var logObject *zerolog.Event

			logType, _ := tengo.ToString(args[0])

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
				return tengo.UndefinedValue, errors.New("unexpected log type")
			}

			m, _ := tengo.ToString(args[1])
			logObject.Msg(m)
			return tengo.UndefinedValue, nil
		}},

		// getEngineSettings()
		// returns the engine settings
		"getEngineSettings": &tengo.UserFunction{Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 0 {
				return tengo.UndefinedValue, errors.New("no arguments can be specified")
			}

			return &e.config, nil
		}},

		// joinPath(path: string...)
		// Joins path names together in an os-specific way
		"joinPath": &tengo.UserFunction{Value: func(args ...tengo.Object) (tengo.Object, error) {
			str := make([]string, len(args))

			for i := range args {
				if args[i].TypeName() != "string" {
					return tengo.UndefinedValue, fmt.Errorf("argument %d must be a string", i+1)
				}

				var ok bool
				str[i], ok = tengo.ToString(args[i])
				if !ok {
					return tengo.UndefinedValue, fmt.Errorf("argument %d could not be converted to a string", i+1)
				}
			}

			return &tengo.String{Value: path.Join(str...)}, nil
		}},

		// sleep(msec: int)
		// sleeps the specified number of milliseconds
		"sleep": &tengo.UserFunction{Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 1 {
				return tengo.UndefinedValue, errors.New("one argument expected")
			}

			val, ok := tengo.ToInt64(args[0])

			if !ok {
				return tengo.UndefinedValue, errors.New("integer expected")
			}

			time.Sleep(time.Millisecond * time.Duration(val))

			return tengo.UndefinedValue, nil
		}},

		// shutdown()
		// shuts down the engine (exits to desktop)
		"shutdown": &tengo.UserFunction{Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 0 {
				return tengo.UndefinedValue, errors.New("no arguments expected")
			}

			log.Info().Msg("engine shutting down due to script SHUTDOWN command")

			e.shutdown = true

			return tengo.UndefinedValue, nil
		}},

		// addLoaderProvider(type: string, path: string)
		// adds a loader to the engine
		"addLoaderProvider": &tengo.UserFunction{Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 2 {
				return tengo.UndefinedValue, errors.New("two arguments expected")
			}

			if args[0].TypeName() != "string" || args[1].TypeName() != "string" {
				return tengo.UndefinedValue, errors.New("both arguments must be string")
			}

			loaderType, _ := tengo.ToString(args[0])
			path, _ := tengo.ToString(args[1])

			switch loaderType {
			case "mpq":
				provider, err := mpqloader.New(path)
				if err != nil {
					return tengo.UndefinedValue, err
				}

				e.loader.AddProvider(provider)
			case "filesystem":
				provider := filesystemloader.New(path)
				e.loader.AddProvider(provider)
			default:
				return tengo.UndefinedValue, fmt.Errorf("unknown loader type: %s", loaderType)
			}

			return tengo.UndefinedValue, nil
		}},

		// loadString(path: string)
		// loads a string from the loader
		"loadString": &tengo.UserFunction{Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 1 {
				return tengo.UndefinedValue, errors.New("expected one parameter")
			}

			val, ok := tengo.ToString(args[0])
			if !ok {
				return tengo.UndefinedValue, errors.New("value must be string")
			}

			file, err := e.loader.Load(val)

			if err != nil {
				return tengo.UndefinedValue, err
			}

			result, err := ioutil.ReadAll(file)
			if err != nil {
				return tengo.UndefinedValue, err
			}

			return &tengo.String{Value: string(result)}, nil
		}},

		// splitString(source: string, splitChars: string)
		// splits a string by the specified split characters
		"splitString": &tengo.UserFunction{Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 2 {
				return tengo.UndefinedValue, errors.New("two parameters expected")
			}

			sourceString, ok := tengo.ToString(args[0])

			if !ok {
				return tengo.UndefinedValue, errors.New("parameters must be of type string")
			}

			splitChars, ok := tengo.ToString(args[1])

			if !ok {
				return tengo.UndefinedValue, errors.New("parameters must be of type string")
			}

			strings := strings.Split(sourceString, splitChars)
			resultArray := make([]tengo.Object, len(strings))

			for i := 0; i < len(strings); i++ {
				resultArray[i] = &tengo.String{Value: strings[i]}
			}

			return &tengo.Array{Value: resultArray}, nil
		}},

		// loadSprite(filePath: string, palette: string) Sprite
		// returns a sprite based on the path and palette
		"loadSprite": &tengo.UserFunction{Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 2 {
				return tengo.UndefinedValue, errors.New("two parameters expected")
			}

			filePath, ok := tengo.ToString(args[0])

			if !ok {
				return tengo.UndefinedValue, errors.New("parameters must be of type string")
			}

			palette, ok := tengo.ToString(args[1])

			if !ok {
				return tengo.UndefinedValue, errors.New("parameters must be of type string")
			}

			result, err := sprite.New(e.loader, filePath, palette)

			if err != nil {
				return nil, err
			}

			return result, nil
		}},

		// exitBootMode()
		// exits boot mode and starts the main rendering system
		"exitBootMode": &tengo.UserFunction{Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 0 {
				return tengo.UndefinedValue, errors.New("no arguments expected")
			}

			log.Info().Msg("Entering game mode")

			e.engineMode = EngineModeGame

			return tengo.UndefinedValue, nil
		}},

		// getRootNode()
		// returns the root entity node for the engine
		"getRootNode": &tengo.UserFunction{Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 0 {
				return tengo.UndefinedValue, errors.New("no arguments expected")
			}

			return e.rootNode, nil
		}},
	})

	script.SetImports(mod)

}
