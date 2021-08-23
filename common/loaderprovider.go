package common

import "io"

type LoaderProvider interface {
	Load(path string) (io.ReadSeekCloser, error)
}
