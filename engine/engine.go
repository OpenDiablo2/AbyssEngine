package engine

import (
	"fmt"
	"io/ioutil"
	"math"
	"runtime"

	lua "github.com/yuin/gopher-lua"

	"github.com/OpenDiablo2/AbyssEngine/common"
	"github.com/OpenDiablo2/AbyssEngine/entity"
	Entity "github.com/OpenDiablo2/AbyssEngine/entity"
	"github.com/OpenDiablo2/AbyssEngine/entity/sprite"
	"github.com/OpenDiablo2/AbyssEngine/loader"
	"github.com/OpenDiablo2/AbyssEngine/loader/filesystemloader"
	"github.com/OpenDiablo2/AbyssEngine/media"
	datPalette "github.com/OpenDiablo2/dat_palette/pkg"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/rs/zerolog/log"
)

// Engine represents the main game engine
type Engine struct {
	config        Configuration
	loader        *loader.Loader
	renderSurface rl.RenderTexture2D
	systemFont    rl.Font
	bootLogo      rl.Texture2D
	bootLoadText  string
	shutdown      bool
	engineMode    EngineMode
	cursorSprite  *sprite.Sprite
	rootNode      *Entity.Entity
	cursorX       int
	cursorY       int
	luaState      *lua.LState
}

func (e *Engine) GetMousePosition() (X, Y int) {
	return e.cursorX, e.cursorY
}

func (e *Engine) GetLanguageCode() string {
	return "eng"
}

// New creates a new instance of the engine
func New(config Configuration) *Engine {
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagVsyncHint)
	rl.InitWindow(800, 600, "Abyss Engine")
	rl.SetTargetFPS(60)
	rl.HideCursor()

	result := &Engine{
		shutdown:      false,
		config:        config,
		engineMode:    EngineModeBoot,
		renderSurface: rl.LoadRenderTexture(800, 600),
		systemFont:    rl.LoadFontFromMemory(".ttf", media.FontDiabloHeavy, int32(len(media.FontDiabloHeavy)), 18, nil, 0),
		rootNode:      entity.New(),
	}

	result.loader = loader.New(result)
	result.loader.AddProvider(filesystemloader.New(config.RootPath))

	logo := rl.LoadImageFromMemory(".png", media.BootLogo, int32(len(media.BootLogo)))
	result.bootLogo = rl.LoadTextureFromImage(logo)
	rl.UnloadImage(logo)

	rl.GenTextureMipmaps(&result.systemFont.Texture)
	rl.SetTextureFilter(result.systemFont.Texture, rl.FilterAnisotropic16x)

	common.PaletteShader = rl.LoadShaderFromMemory(media.StandardVertexShader, media.PaletteFragmentShader)
	common.PaletteShaderLoc = rl.GetShaderLocation(common.PaletteShader, "palette")
	return result
}

// Destroy finalizes the instance of the engine
func (e *Engine) Destroy() {
	rl.UnloadTexture(e.bootLogo)
	rl.UnloadFont(e.systemFont)
}

// Run runs the engine
func (e *Engine) Run() {
	e.bootstrapScripts()

	for !rl.WindowShouldClose() {
		if e.shutdown {
			break
		}

		rl.BeginTextureMode(e.renderSurface)
		rl.ClearBackground(rl.Black)

		switch e.engineMode {
		case EngineModeBoot:
			e.showBootSplash()
		case EngineModeGame:
			e.showGame()
		}

		rl.EndTextureMode()
		e.drawMainSurface()

		if e.engineMode == EngineModeGame {
			e.updateGame()
		}
	}

	e.luaState.Close()
	rl.CloseWindow()
}

func (e *Engine) showGame() {
	e.rootNode.Render()
	if e.cursorSprite != nil {
		e.cursorSprite.Render()
	}
}

func (e *Engine) updateGame() {
	e.rootNode.Update()
	if e.cursorSprite != nil {
		scale := float32(math.Min(float64(rl.GetScreenWidth())/800.0, float64(rl.GetScreenHeight())/600.0))
		xOrigin := (float32(rl.GetScreenWidth()) - (800.0 * scale)) * 0.5
		yOrigin := (float32(rl.GetScreenHeight()) - (600.0 * scale)) * 0.5

		e.cursorX = int((float32(rl.GetMouseX()) - xOrigin) * (1.0 / scale))
		e.cursorSprite.X = e.cursorX

		e.cursorY = int((float32(rl.GetMouseY()) - yOrigin) * (1.0 / scale))
		e.cursorSprite.Y = e.cursorY

		e.cursorSprite.Update()
	}
}

func (e *Engine) showBootSplash() {
	rl.DrawTexture(e.bootLogo, int32(rl.GetScreenWidth()/3)-(e.bootLogo.Width/2),
		int32(rl.GetScreenHeight()/2)-(e.bootLogo.Height/2), rl.White)

	textX := float32(rl.GetScreenWidth()) / 2
	textY := float32(rl.GetScreenHeight()/2) - 20

	rl.DrawTextEx(e.systemFont, "Abyss Engine", rl.Vector2{X: textX, Y: textY}, 18, 0, rl.White)
	rl.DrawTextEx(e.systemFont, "Local Build", rl.Vector2{X: textX, Y: textY + 16}, 18, 0, rl.Gray)
	rl.DrawTextEx(e.systemFont, e.bootLoadText,
		rl.Vector2{X: float32(rl.GetScreenWidth() / 4), Y: float32(rl.GetScreenWidth()/4) * 2.5}, 18, 0, rl.Beige)
}

func (e *Engine) drawMainSurface() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	scale := float32(math.Min(float64(rl.GetScreenWidth())/800.0, float64(rl.GetScreenHeight())/600.0))

	rl.DrawTexturePro(e.renderSurface.Texture,
		rl.Rectangle{Width: float32(e.renderSurface.Texture.Width), Height: float32(-e.renderSurface.Texture.Height)},
		rl.Rectangle{
			X:      (float32(rl.GetScreenWidth()) - (800.0 * scale)) * 0.5,
			Y:      (float32(rl.GetScreenHeight()) - (600.0 * scale)) * 0.5,
			Width:  800.0 * scale,
			Height: 600.0 * scale},
		rl.Vector2{}, 0.0, rl.White)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	rl.DrawTextEx(e.systemFont, fmt.Sprintf("FPS: %d", int(rl.GetFPS())), rl.Vector2{X: 5, Y: 5}, 18, 0, rl.White)
	rl.DrawTextEx(e.systemFont, fmt.Sprintf("GC: %d (%%%d)", int(memStats.NumGC), int(memStats.GCCPUFraction*100)), rl.Vector2{X: 5, Y: 21}, 18, 0, rl.White)
	rl.DrawTextEx(e.systemFont, fmt.Sprintf("Alloc: %d/%d - %0.2fMB", int(memStats.Alloc/1024/1024), int(memStats.TotalAlloc/1024/1024), float32(memStats.Sys/1024/1024)), rl.Vector2{X: 5, Y: 37}, 18, 0, rl.White)

	rl.EndDrawing()
}

func (e *Engine) loadPalette(name string, path string) error {
	if common.PaletteTexture == nil {
		common.PaletteTexture = make(map[string]*common.PalTex)
	}

	println(name, path)

	paletteStream, err := e.loader.Load(path)

	if err != nil {
		return err
	}
	paletteBytes, err := ioutil.ReadAll(paletteStream)

	if err != nil {
		return err
	}

	paletteData, err := datPalette.Decode(paletteBytes)

	if err != nil {
		return err
	}

	colors := make([]byte, 256*4)

	for i := 0; i < 256; i++ {
		if i >= len(paletteData) {
			break
		}

		offset := i * 4
		r, g, b, _ := paletteData[i].RGBA()
		colors[offset] = uint8(r >> 8)
		colors[offset+1] = uint8(g >> 8)
		colors[offset+2] = uint8(b >> 8)
		colors[offset+3] = 255
	}

	colors[3] = 0

	tex := &common.PalTex{}

	tex.Data = colors
	tex.Init = false

	common.PaletteTexture[name] = tex

	return nil
}

func (e *Engine) panic(msg string) {
	// TODO: This should be a UI screen
	log.Fatal().Msg(msg)
	rl.CloseWindow()
}
