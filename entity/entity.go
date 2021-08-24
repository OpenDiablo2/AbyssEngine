package Entity

import (
	"github.com/segmentio/ksuid"
)

type Entity struct {
	id             ksuid.KSUID
	parent         *Entity
	children       []*Entity
	Active         bool
	Visible        bool
	RenderCallback func()
	UpdateCallback func()
}

func New() *Entity {
	result := &Entity{
		id:       ksuid.New(),
		parent:   nil,
		children: make([]*Entity, 0),
		Active:   true,
		Visible:  true,
	}

	return result
}

func (e *Entity) AddChild(entity *Entity) {
	e.children = append(e.children, entity)
}

func (e *Entity) FindChild(id ksuid.KSUID) *Entity {
	// First try a high-level search of direct children
	for idx := range e.children {
		if e.children[idx].id == id {
			return e.children[idx]
		}
	}

	// Not found, do a deep search
	for idx := range e.children {
		result := e.children[idx].FindChild(id)

		if result != nil {
			return result
		}
	}

	// Nothing found..
	return nil
}

func (e *Entity) Render() {
	if !e.Visible || !e.Active {
		return
	}

	if e.RenderCallback != nil {
		e.RenderCallback()
	}

	for idx := range e.children {
		if !e.children[idx].Active || !e.children[idx].Visible {
			continue
		}

		e.children[idx].Render()
	}
}

func (e *Entity) Update() {
	if !e.Active {
		return
	}

	if e.UpdateCallback != nil {
		e.UpdateCallback()
	}

	for idx := range e.children {
		if !e.children[idx].Active {
			continue
		}

		e.children[idx].Update()
	}
}
