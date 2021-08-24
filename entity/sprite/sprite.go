package sprite

import (
	"errors"
	"image/color"
	"io/ioutil"
	"path"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	Entity "github.com/OpenDiablo2/AbyssEngine/entity"

	"github.com/OpenDiablo2/AbyssEngine/common"
	datPalette "github.com/OpenDiablo2/dat_palette/pkg"
	dc6 "github.com/OpenDiablo2/dc6/pkg"
)

type Sprite struct {
	*Entity.Entity

	Sequences       []*dc6.Direction
	palette         []float32
	CurrentSequence int
	CurrentFrame    int
	X               int
	Y               int
	initialized     bool
	Visible         bool
	CellSizeX       int
	CellSizeY       int
	texture         rl.RenderTexture2D
}

func New(loaderProvider common.LoaderProvider, filePath, palette string) (*Sprite, error) {
	result := &Sprite{
		Entity:          Entity.New(),
		X:               0,
		Y:               0,
		initialized:     false,
		Visible:         true,
		CurrentSequence: 0,
		CurrentFrame:    0,
		CellSizeX:       1,
		CellSizeY:       1,
		palette:         make([]float32, 256*3),
	}

	result.RenderCallback = func() { result.render() }
	result.UpdateCallback = func() { result.update() }

	fileExt := strings.ToLower(path.Ext(filePath))

	paletteStream, err := loaderProvider.Load(palette)

	if err != nil {
		return nil, err
	}

	paletteBytes, err := ioutil.ReadAll(paletteStream)

	if err != nil {
		return nil, err
	}

	paletteData, err := datPalette.Decode(paletteBytes)

	if err != nil {
		return nil, err
	}

	fileStream, err := loaderProvider.Load(filePath)

	if err != nil {
		return nil, err
	}

	switch fileExt {
	case ".dcc":
		// TODO: add this
		return nil, errors.New("unsupported file format")
	case ".dc6":
		bytes, err := ioutil.ReadAll(fileStream)

		if err != nil {
			return nil, err
		}

		dc6Res, err := dc6.FromBytes(bytes)

		if err != nil {
			return nil, err
		}

		dc6Res.SetPalette(color.Palette(paletteData))
		result.setPalette(paletteData)
		result.Sequences = dc6Res.Directions

	default:
		return nil, errors.New("unsupported file format")
	}

	_ = fileStream.Close()

	return result, nil
}

func (s *Sprite) setPalette(paletteData datPalette.DAT) {

	for i := 0; i < 256; i++ {
		if i >= len(paletteData) {
			break
		}

		offset := i * 3
		r, g, b, _ := paletteData[i].RGBA()
		s.palette[offset] = float32(r)
		s.palette[offset+1] = float32(g)
		s.palette[offset+2] = float32(b)
	}
}

func (s *Sprite) render() {
	if !s.initialized {
		return
	}
	//rl.DrawRectangle(int32(s.X), int32(s.Y), int32(s.FrameWidth()), int32(s.FrameHeight()), rl.White)
	//rl.SetShaderValueV(common.PaletteShader, common.PaletteShaderLoc, s.palette, rl.ShaderUniformIvec3, 256)
	//rl.BeginShaderMode(common.PaletteShader)
	rl.DrawTexture(s.texture.Texture, int32(s.X), int32(s.Y), rl.White)
	//rl.EndShaderMode()
}

func (s *Sprite) update() {
	if !s.initialized {
		s.initialized = true
		s.initializeTexture()
	}
}

func (s *Sprite) initializeTexture() {
	width := 0
	height := 0

	for i := 0; i < s.CellSizeX; i++ {
		width += int(s.Sequences[s.CurrentSequence].Frames[s.CurrentFrame+i].Width)
	}

	for i := 0; i < s.CellSizeY; i++ {
		height += int(s.Sequences[s.CurrentSequence].Frames[s.CurrentFrame+(i*s.CellSizeX)].Height)
	}

	s.texture = rl.LoadRenderTexture(int32(width), int32(height))
	pixels := make([]rl.Color, width*height)

	targetStartX := 0
	targetStartY := 0

	for cellOffsetY := 0; cellOffsetY < s.CellSizeY; cellOffsetY++ {
		for cellOffsetX := 0; cellOffsetX < s.CellSizeX; cellOffsetX++ {
			cellIndex := s.CurrentFrame + (cellOffsetX + (cellOffsetY * s.CellSizeX))

			frameWidth := s.Sequences[s.CurrentSequence].Frames[cellIndex].Width
			frameHeight := s.Sequences[s.CurrentSequence].Frames[cellIndex].Height

			for y := 0; y < int(frameHeight); y++ {
				idx := targetStartX + ((targetStartY + y) * width)
				for x := 0; x < int(frameWidth); x++ {
					r, g, b, a := s.Sequences[s.CurrentSequence].Frames[cellIndex].At(x, y).RGBA()

					pixels[idx].R = uint8(r)
					pixels[idx].G = uint8(g)
					pixels[idx].B = uint8(b)
					pixels[idx].A = uint8(a)
					idx++
				}
			}

			targetStartX += int(frameWidth)
		}

		targetStartX = 0
		targetStartY += int(s.Sequences[s.CurrentSequence].Frames[(cellOffsetY * s.CellSizeX)].Height)
	}

	rl.UpdateTexture(s.texture.Texture, pixels)

}
