package label

import (
	"errors"
	"io/ioutil"

	"github.com/OpenDiablo2/AbyssEngine/common"
	"github.com/OpenDiablo2/AbyssEngine/node"
	dc6 "github.com/OpenDiablo2/dc6/pkg"
	tblfont "github.com/OpenDiablo2/tbl_font/pkg"
)

type Label struct {
	*node.Node

	FontTable *tblfont.FontTable
	FontGfx   common.SequenceProvider
	Palette   string
}

func New(loaderProvider common.LoaderProvider, fontPath, palette string) (*Label, error) {
	result := &Label{}

	_, ok := common.PaletteTexture[palette]
	if !ok {
		return nil, errors.New("sprite loaded with non-existent palette")
	}
	result.Palette = palette

	fontTableStream, err := loaderProvider.Load(fontPath + ".tbl")
	defer fontTableStream.Close()

	if err != nil {
		return nil, err
	}

	fontData, err := tblfont.Load(fontTableStream)

	if err != nil {
		return nil, err
	}

	result.FontTable = fontData

	fontSpriteStream, err := loaderProvider.Load(fontPath + ".dc6")
	defer fontSpriteStream.Close()

	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(fontSpriteStream)

	if err != nil {
		return nil, err
	}

	spriteData, err := dc6.FromBytes(bytes)

	if err != nil {
		return nil, err
	}

	result.FontGfx = &common.DC6SequenceProvider{Sequences: spriteData.Directions}

	result.RenderCallback = result.render
	result.UpdateCallback = result.update

	return result, nil
}

func (l *Label) render() {

}

func (l *Label) update() {

}
