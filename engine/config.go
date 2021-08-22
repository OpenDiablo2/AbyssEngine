package engine

import (
	"reflect"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"
)

type Configuration struct {
	RootPath     string   `json:"-"`
	MpqLoadOrder []string `json:"mpqLoadOrder"`
}

func (c *Configuration) TypeName() string {
	return "enginesettings"
}

func (c *Configuration) String() string {
	return ""
}

func (c *Configuration) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	panic("implement me")
}

func (c *Configuration) IsFalsy() bool {
	return false
}

func (c *Configuration) Equals(another tengo.Object) bool {
	return false
}

func (c *Configuration) Copy() tengo.Object {
	return nil
}

func (c *Configuration) IndexGet(index tengo.Object) (value tengo.Object, err error) {
	str, _ := tengo.ToString(index)
	r := reflect.ValueOf(c)
	f := reflect.Indirect(r).FieldByName(str)

	switch f.Interface().(type) {
	case string:
		return &tengo.String{Value: f.String()}, nil
	case int:
		return &tengo.Int{Value: f.Int()}, nil
	case []string:
		values := make([]tengo.Object, f.Len())

		for i := 0; i < f.Len(); i++ {
			indexVal := f.Index(i)
			switch indexVal.Interface().(type) {
			case string:
				values[i] = &tengo.String{Value: indexVal.String()}
			case int:
				values[i] = &tengo.Int{Value: indexVal.Int()}
			default:
				panic("todo: types")
			}
		}

		return &tengo.Array{Value: values}, nil
	case bool:
		if f.Bool() {
			return tengo.TrueValue, nil
		}

		return tengo.FalseValue, nil
	default:
		panic("todo: types")
	}

	return tengo.UndefinedValue, nil
}

func (c *Configuration) IndexSet(index, value tengo.Object) error {
	panic("implement me")
}

func (c *Configuration) Iterate() tengo.Iterator {
	panic("implement me")
}

func (c *Configuration) CanIterate() bool {
	return false
}

func (c Configuration) Call(args ...tengo.Object) (ret tengo.Object, err error) {
	panic("implement me")
}

func (c Configuration) CanCall() bool {
	return false
}
