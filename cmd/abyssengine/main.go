package main

import (
	"os"
	"strings"

	"github.com/OpenDiablo2/AbyssEngine/internal/engine"
	"github.com/OpenDiablo2/AbyssEngine/internal/engine/backends/graphicsbackend"
	"github.com/OpenDiablo2/AbyssEngine/internal/engine/backends/graphicsbackend/sdl2graphicsbackend"
	"github.com/OpenDiablo2/AbyssEngine/internal/engine/backends/inputbackend"
	"github.com/OpenDiablo2/AbyssEngine/internal/engine/backends/inputbackend/sdl2inputbackend"
	"github.com/OpenDiablo2/AbyssEngine/internal/engine/configuration"
	"github.com/OpenDiablo2/AbyssEngine/internal/engine/scenemanager"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

const sdl2BackendName = "sdl2"

func main() {
	configureLogging()
	log.Printf("AbyssEngine - The open source ARPG engine")

	fx.New(
		fx.Options(fx.NopLogger),
		fx.Provide(
			// Standard instantiations
			configuration.New,
			engine.New,
			scenemanager.New,

			// Implementation-specific instantiations
			getGraphicsBackend,
			getInputBackend,
		),
		fx.Invoke(run),
	).Run()
}

func run(e *engine.Engine) error {
	return e.Run()
}

func configureLogging() {
	formatter := &log.TextFormatter{
		PadLevelText:     true,
		DisableTimestamp: true,
	}

	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func getGraphicsBackend(config *configuration.Configuration) graphicsbackend.Interface {
	switch strings.ToLower(config.Backend) {
	case sdl2BackendName:
		result, err := sdl2graphicsbackend.Create()
		if err != nil {
			log.Panic(err)
		}

		return result
	default:
		panic("unknown backend")
	}
}

func getInputBackend(config *configuration.Configuration) inputbackend.Interface {
	switch strings.ToLower(config.Backend) {
	case sdl2BackendName:
		result, err := sdl2inputbackend.Create()
		if err != nil {
			log.Panic(err)
		}

		return result
	default:
		panic("unknown backend")
	}
}
