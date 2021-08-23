package sprite

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"

	"github.com/OpenDiablo2/AbyssEngine/common"
	datPalette "github.com/OpenDiablo2/dat_palette/pkg"
	dc6 "github.com/OpenDiablo2/dc6/pkg"
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"
)

type Sprite struct {
	Sequences       []*dc6.Direction
	Palette         datPalette.DAT
	CurrentSequence int
	CurrentFrame    int
	X               int
	Y               int
	Visible         bool
}

func (s *Sprite) TypeName() string {
	return "sprite"
}

func (s *Sprite) String() string {
	return "sprite"
}

func (s *Sprite) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	panic("implement me")
}

func (s *Sprite) IsFalsy() bool {
	panic("implement me")
}

func (s *Sprite) Equals(another tengo.Object) bool {
	panic("implement me")
}

func (s *Sprite) Copy() tengo.Object {
	panic("implement me")
}

func (s *Sprite) IndexGet(index tengo.Object) (value tengo.Object, err error) {
	panic("implement me")
}

func (s *Sprite) IndexSet(index, value tengo.Object) error {
	panic("implement me")
}

func (s *Sprite) Iterate() tengo.Iterator {
	panic("implement me")
}

func (s *Sprite) CanIterate() bool {
	panic("implement me")
}

func (s *Sprite) Call(args ...tengo.Object) (ret tengo.Object, err error) {
	panic("implement me")
}

func (s *Sprite) CanCall() bool {
	panic("implement me")
}

func New(loaderProvider common.LoaderProvider, filePath, palette string) (*Sprite, error) {
	result := &Sprite{
		X:               0,
		Y:               0,
		Visible:         true,
		CurrentSequence: 0,
		CurrentFrame:    0,
	}

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

		result.Palette = paletteData
		result.Sequences = dc6Res.Directions

	default:
		return nil, errors.New("unsupported file format")
	}

	fileStream.Close()

	return result, nil
}
