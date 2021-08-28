package node

import (
	"errors"

	"github.com/segmentio/ksuid"
)

type Node struct {
	Id               ksuid.KSUID
	Parent           *Node
	Children         []*Node
	ChildrenToRemove []*Node
	Active           bool
	Visible          bool
	X                int
	Y                int
	RenderCallback   func()
	UpdateCallback   func()
}

func New() *Node {
	result := &Node{
		Id:       ksuid.New(),
		Parent:   nil,
		Children: make([]*Node, 0),
		Active:   true,
		Visible:  true,
		X:        0,
		Y:        0,
	}

	return result
}

func (e *Node) GetPosition() (X, Y int) {
	if e.Parent == nil {
		return e.X, e.Y
	}

	x, y := e.Parent.GetPosition()

	return e.X + x, e.Y + y
}

func (e *Node) AddChild(entity *Node) error {
	if entity.Parent != nil {
		return errors.New("node already has a Parent")
	}

	e.Children = append(e.Children, entity)
	entity.Parent = e

	return nil
}

func (e *Node) RemoveAllChildren() {
	for idx := range e.Children {
		e.Children[idx].Parent = nil
	}

	e.Children = make([]*Node, 0)

	return
}

func (e *Node) RemoveChild(node *Node) {
	for idx := range e.Children {
		if e.Children[idx] != node {
			continue
		}

		e.Children[idx].Parent = nil
		newChildren := e.Children[:idx]
		newChildren = append(newChildren, e.Children[idx+1:]...)

		e.Children = newChildren

		return
	}
}

func (e *Node) FindChild(id ksuid.KSUID) *Node {
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

func (e *Node) Render() {
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

func (e *Node) Update() {
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
