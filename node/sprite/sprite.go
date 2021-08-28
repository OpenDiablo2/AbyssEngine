package sprite

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"

	"github.com/OpenDiablo2/AbyssEngine/common"
	"github.com/OpenDiablo2/AbyssEngine/node"
	dc6 "github.com/OpenDiablo2/dc6/pkg"
	dcc "github.com/OpenDiablo2/dcc/pkg"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Sprite struct {
	*node.Node

	mousePosProvider  common.MousePositionProvider
	Sequences         common.SequenceProvider
	palette           string
	CurrentSequence   int
	CurrentFrame      int
	initialized       bool
	Visible           bool
	CellSizeX         int
	CellSizeY         int
	isPressed         bool
	texture           rl.Texture2D
	onMouseButtonDown func()
	onMouseButtonUp   func()
}

func New(loaderProvider common.LoaderProvider, mousePosProvider common.MousePositionProvider,
	filePath, palette string) (*Sprite, error) {
	result := &Sprite{
		Node:             node.New(),
		mousePosProvider: mousePosProvider,
		initialized:      false,
		Visible:          true,
		CurrentSequence:  0,
		CurrentFrame:     0,
		CellSizeX:        1,
		CellSizeY:        1,
		isPressed:        false,
		palette:          palette,
	}

	result.RenderCallback = func() { result.render() }
	result.UpdateCallback = func() { result.update() }

	fileExt := strings.ToLower(path.Ext(filePath))

	fileStream, err := loaderProvider.Load(filePath)

	if err != nil {
		return nil, err
	}

	_, ok := common.PaletteTexture[palette]
	if !ok {
		return nil, errors.New("Sprite loaded with non-existent palette")
	}

	switch fileExt {
	case ".dcc":
		bytes, err := ioutil.ReadAll(fileStream)

		if err != nil {
			return nil, err
		}

		dccRes, err := dcc.FromBytes(bytes)

		if err != nil {
			return nil, err
		}

		result.Sequences = &DCCSequenceProvider{dccRes.Directions()}

	case ".dc6":
		bytes, err := ioutil.ReadAll(fileStream)

		if err != nil {
			return nil, err
		}

		dc6Res, err := dc6.FromBytes(bytes)

		if err != nil {
			return nil, err
		}

		result.Sequences = &DC6SequenceProvider{Sequences: dc6Res.Directions}

	default:
		return nil, errors.New("unsupported file format")
	}

	_ = fileStream.Close()

	return result, nil
}

func (s *Sprite) setPalette(palette string) {
	s.palette = palette
}

func (s *Sprite) render() {
	if !s.initialized || !s.Visible || !s.Active {
		return
	}

	tex := common.PaletteTexture[s.palette]
	if !tex.Init {
		img := rl.NewImage(tex.Data, 256, 1, 1, rl.UncompressedR8g8b8a8)
		tex.Texture = rl.LoadTextureFromImage(img)

		tex.Init = true
	}

	rl.BeginShaderMode(common.PaletteShader)
	rl.SetShaderValueTexture(common.PaletteShader, common.PaletteShaderLoc, tex.Texture)
	rl.DrawTexture(s.texture, int32(s.X), int32(s.Y), rl.White)
	rl.EndShaderMode()
}

func (s *Sprite) update() {
	if !s.initialized {
		s.initialized = true
		s.initializeTexture()
	}

	if !s.isPressed && rl.IsMouseButtonDown(rl.MouseLeftButton) {

		mx, my := s.mousePosProvider.GetMousePosition()
		posX, posY := s.GetPosition()

		if mx < posX || my < posY || mx >= (posX+int(s.texture.Width)) || my >= (posY+int(s.texture.Height)) {
			return
		}

		s.isPressed = true

		if s.onMouseButtonDown != nil {
			s.onMouseButtonDown()
		}
	} else {
		s.isPressed = false
	}
}

func (s *Sprite) initializeTexture() {
	width := 0
	height := 0

	for i := 0; i < s.CellSizeX; i++ {
		width += s.Sequences.FrameWidth(s.CurrentSequence, s.CurrentFrame+i)
	}

	for i := 0; i < s.CellSizeY; i++ {
		height += s.Sequences.FrameHeight(s.CurrentSequence, s.CurrentFrame+(i*s.CellSizeX))
	}

	pixels := make([]byte, width*height)

	targetStartX := 0
	targetStartY := 0

	for cellOffsetY := 0; cellOffsetY < s.CellSizeY; cellOffsetY++ {
		for cellOffsetX := 0; cellOffsetX < s.CellSizeX; cellOffsetX++ {
			cellIndex := s.CurrentFrame + (cellOffsetX + (cellOffsetY * s.CellSizeX))

			frameWidth := s.Sequences.FrameWidth(s.CurrentSequence, cellIndex)
			frameHeight := s.Sequences.FrameHeight(s.CurrentSequence, cellIndex)

			for y := 0; y < frameHeight; y++ {
				idx := targetStartX + ((targetStartY + y) * width)
				for x := 0; x < frameWidth; x++ {
					c := s.Sequences.GetColorIndexAt(s.CurrentSequence, cellIndex, x, y)

					pixels[idx] = c
					idx++
				}
			}

			targetStartX += int(frameWidth)
		}

		targetStartX = 0
		targetStartY += s.Sequences.FrameHeight(s.CurrentSequence, cellOffsetY*s.CellSizeX)
	}

	img := rl.NewImage(pixels, int32(width), int32(height), 1, rl.UncompressedGrayscale)
	s.texture = rl.LoadTextureFromImage(img)

}
