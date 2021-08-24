package sprite

import (
	"errors"
	"fmt"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (s *Sprite) TypeName() string {
	return "entity"
}

func (s *Sprite) String() string {
	return "entity'"
}

func (s *Sprite) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	panic("implement me")
}

func (s *Sprite) IsFalsy() bool {
	return false
}

func (s *Sprite) Equals(another tengo.Object) bool {
	panic("implement me")
}

func (s *Sprite) Copy() tengo.Object {
	panic("implement me")
}

func (s *Sprite) IndexGet(index tengo.Object) (value tengo.Object, err error) {
	indexStr, ok := tengo.ToString(index)

	if !ok {
		return nil, errors.New("invalid index")
	}

	switch indexStr {
	case "x":
		return &tengo.Int{Value: int64(s.X)}, nil
	case "y":
		return &tengo.Int{Value: int64(s.Y)}, nil
	case "node":
		return s.Entity, nil
	case "setPosition":
		return &tengo.UserFunction{
			Name: "appendChild",
			Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
				if len(args) != 2 {
					return nil, errors.New("expected two arguments")
				}

				posX, ok := tengo.ToInt(args[0])

				if !ok {
					return nil, errors.New("first argument must be int")
				}

				posY, ok := tengo.ToInt(args[0])

				if !ok {
					return nil, errors.New("first argument must be int")
				}

				s.X = posX
				s.Y = posY

				return s, nil
			},
		}, nil
	case "setCellSize":
		return &tengo.UserFunction{
			Name: "setCellSize",
			Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
				if len(args) != 2 {
					return nil, errors.New("expected two arguments")
				}

				sizeX, ok := tengo.ToInt(args[0])

				if !ok {
					return nil, errors.New("first argument must be int")
				}

				sizeY, ok := tengo.ToInt(args[1])

				if !ok {
					return nil, errors.New("first argument must be int")
				}

				s.CellSizeX = sizeX
				s.CellSizeY = sizeY
				s.initialized = false
				rl.UnloadRenderTexture(s.texture)

				return s, nil
			},
		}, nil
	case "setActive":
		return &tengo.UserFunction{
			Name: "setActive",
			Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
				if len(args) != 1 {
					return nil, errors.New("expected one argument")
				}

				value, ok := tengo.ToBool(args[0])

				if !ok {
					return nil, errors.New("first argument must be boolean")
				}

				s.Active = value

				return s, nil
			},
		}, nil
	case "setVisible":
		return &tengo.UserFunction{
			Name: "setVisible",
			Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
				if len(args) != 1 {
					return nil, errors.New("expected one argument")
				}

				value, ok := tengo.ToBool(args[0])

				if !ok {
					return nil, errors.New("first argument must be boolean")
				}

				s.Visible = value

				return s, nil
			},
		}, nil
	}

	return nil, fmt.Errorf("invalid index: %s", indexStr)

}

func (s *Sprite) IndexSet(index, value tengo.Object) error {
	indexStr, ok := tengo.ToString(index)

	if !ok {
		return errors.New("invalid index")
	}

	return fmt.Errorf("invalid index: %s", indexStr)
}

func (s *Sprite) Iterate() tengo.Iterator {
	panic("implement me")
}

func (s *Sprite) CanIterate() bool {
	return false
}

func (s *Sprite) Call(args ...tengo.Object) (ret tengo.Object, err error) {
	panic("implement me")
}

func (s *Sprite) CanCall() bool {
	panic("implement me")
}
