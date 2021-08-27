package common

import rl "github.com/gen2brain/raylib-go/raylib"

type PalTex struct {
	Texture rl.Texture2D
	Data []byte
	Init bool
}

//TODO: Yeah yeah, move this out
var (
	PaletteShader    rl.Shader
	PaletteShaderLoc int32
	PaletteTexture   map[string]*PalTex
)
