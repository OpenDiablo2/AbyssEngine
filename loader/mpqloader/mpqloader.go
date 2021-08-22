package mpqloader

import (
	"fmt"
	"io"
	"strings"

	"github.com/OpenDiablo2/mpq"
)

type MpqLoader struct {
	mpq *mpq.MPQ
}

func (m *MpqLoader) Name() string {
	return "MPQ Loader"
}

func (m *MpqLoader) Exists(path string) bool {
	if len(path) == 0 {
		return false
	}

	path = strings.ReplaceAll(path, "/", "\\")

	return m.mpq.Contains(path)
}

func (m *MpqLoader) Load(path string) (io.ReadSeekCloser, error) {
	path = strings.ReplaceAll(path, "/", "\\")

	if !m.Exists(path) {
		return nil, fmt.Errorf("could not locate file \"%s\" in \"%s\"", path, m.mpq.Path())
	}

	return m.mpq.ReadFileStream(path)
}

func New(fileName string) (*MpqLoader, error) {
	result := &MpqLoader{}

	mpq, err := mpq.FromFile(fileName)

	if err != nil {
		return nil, err
	}

	result.mpq = mpq

	return result, nil
}
