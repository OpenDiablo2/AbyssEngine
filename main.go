package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path"

	"github.com/OpenDiablo2/AbyssEngine/engine"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var runPath string

func initFlags() {
	flag.StringVar(&runPath, "path", "", "path to the engine runtime files")
	flag.Parse()

	if runPath == "" {
		runPath, _ = os.Getwd()
	}
}

func initLogging() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {
	initFlags()
	initLogging()

	log.Info().Msg("Abyss Engine")
	log.Debug().Msgf("Runtime Path: %s", runPath)

	rl.SetTraceLogCallback(func(logLevel int, s string) {
		[]func() *zerolog.Event{
			log.Trace,
			log.Debug,
			log.Info,
			log.Warn,
			log.Error,
			log.Fatal,
		}[logLevel-1]().Msg(s)
	})

	engineConfig := engine.Configuration{
		RootPath: runPath,
	}

	jsonFile, err := os.Open(path.Join(runPath, "config.json"))
	if err == nil {
		bytes, _ := ioutil.ReadAll(jsonFile)
		_ = json.Unmarshal(bytes, &engineConfig)
		_ = jsonFile.Close()
	}

	coreEngine := engine.New(engineConfig)

	coreEngine.Run()
	coreEngine.Destroy()

}
