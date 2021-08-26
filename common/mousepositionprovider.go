package common

type MousePositionProvider interface {
	GetMousePosition() (X, Y int)
}
