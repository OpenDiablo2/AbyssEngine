package loader

import "io"

type LoaderProvider interface {
	Name() string
	Exists(path string) bool
	Load(path string) (io.ReadSeekCloser, error)
}

