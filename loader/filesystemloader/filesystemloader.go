package filesystemloader

import (
	"io"
	"os"
	"path"
)

type FileSystemLoader struct {
	basePath string
}

func (f *FileSystemLoader) Name() string {
	return "FileSystem Loader"
}

func (f *FileSystemLoader) Exists(p string) bool {
	_, err := os.Stat(path.Join(f.basePath, p))

	return !os.IsNotExist(err)
}

func (f *FileSystemLoader) Load(p string) (io.ReadSeekCloser, error) {
	return os.Open(path.Join(f.basePath, p))
}

func New(basePath string) *FileSystemLoader {
	result := &FileSystemLoader{
		basePath: basePath,
	}

	return result
}
