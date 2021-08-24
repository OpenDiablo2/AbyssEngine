package Entity

import (
	"errors"
	"fmt"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"
)

func (e *Entity) TypeName() string {
	return "entity"
}

func (e *Entity) String() string {
	return "entity'"
}

func (e *Entity) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	panic("implement me")
}

func (e *Entity) IsFalsy() bool {
	return false
}

func (e *Entity) Equals(another tengo.Object) bool {
	panic("implement me")
}

func (e *Entity) Copy() tengo.Object {
	panic("implement me")
}

func (e *Entity) IndexGet(index tengo.Object) (value tengo.Object, err error) {
	indexStr, ok := tengo.ToString(index)

	if !ok {
		return nil, errors.New("invalid index")
	}

	switch indexStr {
	case "id":
		return &tengo.String{Value: e.id.String()}, nil
	case "appendChild":
		return &tengo.UserFunction{
			Name: "appendChild",
			Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
				if len(args) != 1 {
					return nil, errors.New("expected one argument")
				}

				target, ok := tengo.ToInterface(args[0]).(*Entity)

				if !ok {
					return nil, errors.New("entity expected")
				}

				e.AddChild(target)

				return e, nil
			},
		}, nil
	}

	return nil, fmt.Errorf("unknown index: %s", indexStr)

}

func (e *Entity) IndexSet(index, value tengo.Object) error {
	panic("implement me")
}

func (e *Entity) Iterate() tengo.Iterator {
	panic("implement me")
}

func (e *Entity) CanIterate() bool {
	return false
}

func (e *Entity) Call(args ...tengo.Object) (ret tengo.Object, err error) {
	panic("implement me")
}

func (e *Entity) CanCall() bool {
	panic("implement me")
}
