package label

import "github.com/OpenDiablo2/AbyssEngine/node"

type Label struct {
	*node.Node
}

func New(filePath, palette string) *Label {
	result := &Label{}
	return result
}
