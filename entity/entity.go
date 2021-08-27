package entity

import (
	"errors"

	"github.com/segmentio/ksuid"
)

type Entity struct {
	Id             ksuid.KSUID
	Parent         *Entity
	Children       []*Entity
	Active         bool
	Visible        bool
	X              int
	Y              int
	RenderCallback func()
	UpdateCallback func()
}

func New() *Entity {
	result := &Entity{
		Id:       ksuid.New(),
		Parent:   nil,
		Children: make([]*Entity, 0),
		Active:   true,
		Visible:  true,
		X:        0,
		Y:        0,
	}

	return result
}

func (e *Entity) GetPosition() (X, Y int) {
	if e.Parent == nil {
		return e.X, e.Y
	}

	x, y := e.Parent.GetPosition()

	return e.X + x, e.Y + y
}

func (e *Entity) AddChild(entity *Entity) error {
	if entity.Parent != nil {
		return errors.New("entity already has a Parent")
	}

	e.Children = append(e.Children, entity)
	entity.Parent = e

	return nil
}

func (e *Entity) FindChild(id ksuid.KSUID) *Entity {
	// First try a high-level search of direct Children
	for idx := range e.Children {
		if e.Children[idx].Id == id {
			return e.Children[idx]
		}
	}

	// Not found, do a deep search
	for idx := range e.Children {
		result := e.Children[idx].FindChild(id)

		if result != nil {
			return result
		}
	}

	// Nothing found...
	return nil
}

func (e *Entity) Render() {
	if !e.Visible || !e.Active {
		return
	}

	if e.RenderCallback != nil {
		e.RenderCallback()
	}

	for idx := range e.Children {
		if !e.Children[idx].Active || !e.Children[idx].Visible {
			continue
		}

		e.Children[idx].Render()
	}
}

func (e *Entity) Update() {
	if !e.Active {
		return
	}

	if e.UpdateCallback != nil {
		e.UpdateCallback()
	}

	for idx := range e.Children {
		if !e.Children[idx].Active {
			continue
		}

		e.Children[idx].Update()
	}
}
