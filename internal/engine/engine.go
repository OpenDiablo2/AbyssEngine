package engine

import (
	"bufio"
	"log"
	"os"
	"time"

	"github.com/OpenDiablo2/AbyssEngine/internal/engine/scenemanager"

	"github.com/OpenDiablo2/AbyssEngine/internal/engine/backends/inputbackend"

	"github.com/OpenDiablo2/AbyssEngine/internal/engine/backends/graphicsbackend"

	"github.com/OpenDiablo2/AbyssEngine/internal/engine/configuration"
	"github.com/gravestench/akara"
)

// Engine represents an instance of the Abyss Engine.
type Engine struct {
	ecs          *akara.World
	gfx          graphicsbackend.Interface
	input        inputbackend.Interface
	sceneManager *scenemanager.SceneManager
}

// New creates a new instance of the abyss engine
func New(config *configuration.Configuration,
	graphicsBackend graphicsbackend.Interface,
	inputBackend inputbackend.Interface,
	sceneManager *scenemanager.SceneManager,
) (*Engine, error) {
	result := &Engine{
		gfx:          graphicsBackend,
		input:        inputBackend,
		sceneManager: sceneManager,
	}

	result.configureECS()

	return result, nil
}

// Run runs the engine
func (engine *Engine) Run() error {
	lastUpdateTime := time.Now()
	go func() {
		engine.handleDebugger()
	}()

	for {
		currentTime := time.Now()
		timeDelta := currentTime.Sub(lastUpdateTime)
		lastUpdateTime = currentTime

		if err := engine.ecs.Update(timeDelta); err != nil {
			return err
		}

		if err := engine.gfx.Render(); err != nil {
			return err
		}

		if err := engine.input.Process(); err != nil {
			return err
		}
	}
}

func (engine *Engine) configureECS() {
	cfg := akara.NewWorldConfig().
		With(&configuration.Configuration{})

	engine.ecs = akara.NewWorld(cfg)
}

func (engine *Engine) handleDebugger() {
	reader := bufio.NewReader(os.Stdin)
	for {
		res, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Couldn't read stdin")
		}

		log.Printf("Recieved: %q", res)
	}
}
